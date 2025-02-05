"""
# Copyright 2021 21CN Corporation Limited
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""

# -*- coding: utf-8 -*-

import json
import os
import threading
import uuid
import zipfile

from heatclient.common import template_utils
from heatclient.exc import HTTPNotFound
from pony.orm import db_session, commit

import utils
from core import openstack_utils
from core.csar.pkg import get_hot_yaml_path, CsarPkg
from core.log import logger
from core.models import AppInsMapper, InstantiateRequest, UploadCfgRequest, UploadPackageRequest, BaseRequest
from internal.lcmservice import lcmservice_pb2_grpc
from internal.lcmservice.lcmservice_pb2 import TerminateResponse, \
    QueryResponse, UploadCfgResponse, \
    RemoveCfgResponse, DeletePackageResponse, UploadPackageResponse, \
    WorkloadEventsResponse

LOG = logger


def start_check_stack_status(app_instance_id):
    """
    start_check_stack_status
    Args:
        app_instance_id:
    """
    thread_timer = threading.Timer(5, check_stack_status, [app_instance_id])
    thread_timer.start()


@db_session
def check_stack_status(app_instance_id):
    """
    check_stack_status
    Args:
        app_instance_id:
    """
    app_ins_mapper = AppInsMapper.get(app_instance_id=app_instance_id)
    if not app_ins_mapper:
        LOG.debug('app ins: %s db record not found', app_instance_id)
        return
    heat = openstack_utils.create_heat_client(app_ins_mapper.host_ip)
    stack_resp = heat.stacks.get(app_ins_mapper.stack_id)
    if stack_resp is None and app_ins_mapper.operational_status == 'Terminating':
        app_ins_mapper.delete()
        LOG.debug('finish terminate app ins %s', app_instance_id)
        return
    if stack_resp.status == 'COMPLETE' or stack_resp.status == 'FAILED':
        LOG.debug('app ins: %s, stack_status: %s, reason: %s',
                  app_instance_id,
                  stack_resp.stack_status,
                  stack_resp.stack_status_reason)
        if stack_resp.action == 'CREATE' and stack_resp.stack_status == 'CREATE_COMPLETE':
            app_ins_mapper.operational_status = utils.INSTANTIATED
            app_ins_mapper.operation_info = stack_resp.stack_status_reason
            LOG.debug('finish instantiate app ins %s', app_instance_id)
        elif stack_resp.action == 'DELETE' and stack_resp.stack_status == 'DELETE_COMPLETE':
            app_ins_mapper.delete()
            LOG.debug('finish terminate app ins %s', app_instance_id)
        else:
            app_ins_mapper.operation_info = stack_resp.stack_status_reason
            app_ins_mapper.operational_status = utils.FAILURE
            LOG.debug('failed action %s app ins %s', stack_resp.action, app_instance_id)
    else:
        LOG.debug('app ins %s status not updated, waite next...')
        start_check_stack_status(app_instance_id)


def validate_input_params(param):
    """
    check_stack_status
    Args:
        param:
    Returns:
        host_ip
    """
    access_token = param.access_token
    host_ip = param.host_ip
    LOG.debug('param hostIp: %s', host_ip)
    LOG.debug('param accessToken: %s', access_token)
    if not utils.validate_access_token(access_token):
        return None
    if not utils.validate_ipv4_address(host_ip):
        return None
    return host_ip


def _get_output_data(output_list, heat, stack_id):
    """
    获取output数据
    """
    response = {
        'code': 200,
        'msg': 'ok',
        'data': []
    }
    for item in output_list['outputs']:
        output = heat.stacks.output_show(stack_id, item['output_key'])
        output_value = output['output']['output_value']
        item = {
            'vmId': output_value['vmId'],
            'vncUrl': output_value['vncUrl'],
            'networks': []
        }
        if 'networks' in output_value:
            for net_name, ip_data in output_value['networks'].items():
                if utils.validate_uuid(net_name):
                    continue
                network = {
                    'name': net_name,
                    'ip': ip_data[0]['addr']
                }
                item['networks'].append(network)
            response['data'].append(item)
    return response


class AppLcmService(lcmservice_pb2_grpc.AppLCMServicer):
    """
    AppLcmService
    """
    def uploadPackage(self, request_iterator, context):
        """
        上传app包
        :param request_iterator:
        :param context:
        :return:
        """
        LOG.info('receive upload package msg...')
        res = UploadPackageResponse(status=utils.FAILURE)

        parameters = UploadPackageRequest(request_iterator)

        host_ip = validate_input_params(parameters)
        if host_ip is None:
            parameters.delete_tmp()
            return res

        app_package_id = parameters.app_package_id
        if app_package_id is None:
            LOG.debug('appPackageId is required')
            parameters.delete_tmp()
            return res
        app_package_path = utils.APP_PACKAGE_DIR + '/' + host_ip + '/' + parameters.app_package_id
        if utils.exists_path(app_package_path):
            LOG.debug('app package exist')
            parameters.delete_tmp()
            return res
        utils.create_dir(app_package_path)
        try:
            LOG.debug('unzip package')
            with zipfile.ZipFile(parameters.tmp_package_file_path) as zip_file:
                namelist = zip_file.namelist()
                for file in namelist:
                    zip_file.extract(file, app_package_path)
            pkg = CsarPkg(app_package_path)
            pkg.translate()
            res.status = utils.SUCCESS
        except Exception as exception:
            LOG.error(exception, exc_info=True)
            utils.delete_dir(app_package_path)
        finally:
            parameters.delete_tmp()
        return res

    def deletePackage(self, request, context):
        """
        删除app包
        :param request:
        :param context:
        :return:
        """
        LOG.info('receive delete package msg...')
        res = DeletePackageResponse(status=utils.FAILURE)

        host_ip = validate_input_params(BaseRequest(request))
        if host_ip is None:
            return res

        app_package_id = request.appPackageId
        if not app_package_id:
            return res

        app_package_path = utils.APP_PACKAGE_DIR + '/' + host_ip + '/' + app_package_id
        utils.delete_dir(app_package_path)

        res.status = utils.SUCCESS
        return res

    @db_session
    def instantiate(self, request, context):
        """
        app 实例化
        :param request:
        :param context:
        :return:
        """
        LOG.info('receive instantiate msg...')
        res = TerminateResponse(status=utils.FAILURE)

        parameter = InstantiateRequest(request)

        LOG.debug('校验access token, host ip')
        host_ip = validate_input_params(parameter)
        if host_ip is None:
            return res

        LOG.debug('获取实例ID')
        app_instance_id = parameter.app_instance_id
        if app_instance_id is None:
            return res

        LOG.debug('查询数据库是否存在相同记录')
        app_ins_mapper = AppInsMapper.get(app_instance_id=app_instance_id)
        if app_ins_mapper is not None:
            LOG.info('app ins %s exist', app_instance_id)
            return res

        LOG.debug('读取包的hot文件')
        hot_yaml_path = get_hot_yaml_path(parameter.app_package_path)
        if hot_yaml_path is None:
            return res

        LOG.debug('构建heat参数')
        tpl_files, template = template_utils.get_template_contents(template_file=hot_yaml_path)
        fields = {
            'stack_name': 'eg-' + ''.join(str(uuid.uuid4()).split('-'))[0:8],
            'template': template,
            'files': dict(list(tpl_files.items()))
        }
        LOG.debug('init heat client')
        heat = openstack_utils.create_heat_client(host_ip)
        try:
            LOG.debug('发送创建stack请求')
            stack_resp = heat.stacks.create(**fields)
        except Exception as exception:
            LOG.error(exception, exc_info=True)
            return res
        AppInsMapper(app_instance_id=app_instance_id,
                     host_ip=host_ip,
                     stack_id=stack_resp['stack']['id'],
                     operational_status=utils.INSTANTIATING)
        LOG.debug('更新数据库')
        commit()

        LOG.debug('开始更新状态定时任务')
        start_check_stack_status(app_instance_id=app_instance_id)

        res.status = utils.SUCCESS
        LOG.debug('消息处理完成')
        return res

    @db_session
    def terminate(self, request, context):
        """
        销毁实例
        :param request:
        :param context:
        :return:
        """
        LOG.info('receive terminate msg...')
        res = TerminateResponse(status=utils.FAILURE)

        LOG.debug('校验token, host ip')
        host_ip = validate_input_params(BaseRequest(request))
        if host_ip is None:
            return res

        LOG.debug('获取实例ID')
        app_instance_id = request.appInstanceId
        if app_instance_id is None:
            return res

        LOG.debug('查询数据库')
        app_ins_mapper = AppInsMapper.get(app_instance_id=app_instance_id)
        if app_ins_mapper is None:
            res.status = utils.SUCCESS
            return res

        LOG.debug('初始化openstack客户端')
        heat = openstack_utils.create_heat_client(host_ip)
        try:
            LOG.debug('发送删除请求')
            heat.stacks.delete(app_ins_mapper.stack_id)
        except HTTPNotFound:
            LOG.debug('stack不存在')
        except Exception as exception:
            LOG.error(exception, exc_info=True)
            return res

        app_ins_mapper.operational_status = utils.TERMINATING
        LOG.debug('更新数据库状态')
        commit()

        LOG.debug('开始状态更新定时任务')
        start_check_stack_status(app_instance_id=app_instance_id)

        res.status = utils.SUCCESS
        LOG.debug('处理请求完成')
        return res

    @db_session
    def query(self, request, context):
        """
        实例信息查询
        :param request:
        :param context:
        :return:
        """
        LOG.info('receive query msg...')
        res = QueryResponse(response='{"code": 500, "msg": "server error"}')

        host_ip = validate_input_params(BaseRequest(request))
        if host_ip is None:
            res.response = '{"code":400}'
            return res

        app_instance_id = request.appInstanceId
        if app_instance_id is None:
            res.response = '{"code":400}'
            return res

        app_ins_mapper = AppInsMapper.get(app_instance_id=app_instance_id)
        if app_ins_mapper is None:
            res.response = '{"code":404}'
            return res

        heat = openstack_utils.create_heat_client(host_ip)
        output_list = heat.stacks.output_list(app_ins_mapper.stack_id)

        response = _get_output_data(output_list, heat, app_ins_mapper.stack_id)

        res.response = json.dumps(response)
        return res

    @db_session
    def workloadEvents(self, request, context):
        """
        工作负载事件查询
        :param request:
        :param context:
        :return:
        """
        LOG.info('receive workload describe msg...')
        res = WorkloadEventsResponse(response='{"code":500}')

        host_ip = validate_input_params(BaseRequest(request))
        if host_ip is None:
            return res

        app_instance_id = request.appInstanceId
        if app_instance_id is None:
            return res

        app_ins_mapper = AppInsMapper.get(app_instance_id=app_instance_id)
        if app_ins_mapper is None:
            LOG.debug('app实例 %s 不存在', app_instance_id)
            res.response = '{"code":404}'
            return res

        heat = openstack_utils.create_heat_client(host_ip)

        events = heat.events.list(stack_id=app_ins_mapper.stack_id)
        vm_describe_info = {}
        for event in events:
            if event.resource_name in vm_describe_info:
                vm_describe_info[event.resource_name]['events'].append({
                    'eventTime': event.event_time,
                    'resourceStatus': event.resource_status,
                    'resourceStatusReason': event.resource_status_reason
                })
            else:
                vm_describe_info[event.resource_name] = {
                    'resourceName': event.resource_name,
                    'logicalResourceId': event.logical_resource_id,
                    'physicalResourceId': event.physical_resource_id,
                    'events': [
                        {
                            'eventTime': event.event_time,
                            'resourceStatus': event.resource_status,
                            'resourceStatusReason': event.resource_status_reason
                        }
                    ]
                }
        response_data = []
        for key, value in vm_describe_info.items():
            response_data.append(value)
        res.response = json.dumps(response_data)
        return res

    def uploadConfig(self, request_iterator, context):
        """
        上传openstack配置文件
        :param request_iterator: 流式传输
        :param context:
        :return:
        """
        LOG.info('receive uploadConfig msg...')
        res = UploadCfgResponse(status=utils.FAILURE)

        parameter = UploadCfgRequest(request_iterator)

        host_ip = validate_input_params(parameter)
        if host_ip is None:
            return res

        config_file = parameter.config_file
        if config_file is None:
            return res

        config_path = utils.RC_FILE_DIR + '/' + host_ip

        try:
            with open(config_path, 'wb') as new_file:
                new_file.write(config_file)
            openstack_utils.set_rc(host_ip)
            openstack_utils.clear_glance_client(host_ip)
            res.status = utils.SUCCESS
        except Exception as exception:
            LOG.error(exception, exc_info=True)

        return res

    def removeConfig(self, request, context):
        """
        删除openstack 配置文件
        :param request: 请求体
        :param context: 上下文信息
        :return: Success/Failure
        """
        LOG.info('receive removeConfig msg...')
        res = RemoveCfgResponse(status=utils.FAILURE)

        host_ip = validate_input_params(BaseRequest(request))
        if not host_ip:
            return res

        config_path = utils.RC_FILE_DIR + '/' + host_ip
        try:
            os.remove(config_path)
            openstack_utils.del_rc(host_ip)
            openstack_utils.clear_glance_client(host_ip)
            res.status = utils.SUCCESS
        except OSError:
            res.status = utils.SUCCESS
        except Exception as exception:
            LOG.error(exception, exc_info=True)
            return res

        LOG.debug('host configuration file deleted successfully')
        return res
