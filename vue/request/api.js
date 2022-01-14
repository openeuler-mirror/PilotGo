// api请求接口文件
import { request } from './request'  //引入axios封装的request请求

const API1 = '/plugin/prometheus/api/v1'
export function getChartName(time) {
  return request({
    url: API1 + '/label/__name__/values?_=' + time,
    method: 'get',
  })
}

export function getChart(url) {
  return request({
    url: API1 + url,
    method: 'get',
  })
}

export function loginByEmail(data) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}
export function logout() {
  return request({
    url: '/user/logout',
    method: 'post',
  })
}

export function getPlugins() {
  return request({
    url: '/plugin',
    method: 'get',
  })
}

export function insertPlugin(data) {
  return request({
    url: '/plugin',
    method: 'post',
    data
  })
}

export function deletePlugins(data) {
  return request({
    url: '/plugin',
    method: 'delete',
    data
  })
}

export function getOverview() {
  return request({
    url: '/overview',
    method: 'get'
  })
}

export function FirewallConfig(data) {
  return request({
    url: '/firewall/config',
    method: 'post',
    data
  })
}

export function FirewallStop(data) {
  return request({
    url: '/firewall/stop',
    method: 'post',
    data
  })
}

export function FirewallRestart(data) {
  return request({
    url: '/firewall/restart',
    method: 'post',
    data
  })
}

export function FirewallReload(data) {
  return request({
    url: '/firewall/reload',
    method: 'post',
    data
  })
}

export function FirewallAddZonePort(data) {
  return request({
    url: '/firewall/addzp',
    method: 'post',
    data
  })
}

export function FirewallDelZonePort(data) {
  return request({
    url: '/firewall/delzp',
    method: 'post',
    data
  })
}

export function FirewallAddZonePortPermanent(data) {
  return request({
    url: '/firewall/addzpp',
    method: 'post',
    data
  })
}
