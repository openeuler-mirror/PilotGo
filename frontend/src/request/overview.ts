import request from './request';

// 获取机器集群概览信息
export function machinesOverview() {
  return request({
    url: '/overview/info',
    method: 'get',
  });
}

// 获取各个部门机器集群概览信息
export function departMachinesOverview() {
  return request({
    url: '/overview/depart_info',
    method: 'get',
  });
}