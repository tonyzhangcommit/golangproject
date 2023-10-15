import request from '@/utils/request'

export function getList(params) {
  return request({
    url: '/vue-admin-template/table/list',
    method: 'get',
    params
  })
}

// 获取用户列表
export function getuserlist(params) {
  return request({
    url: '/superadmin/users',
    method: 'get',
    params
  })
}
// 获取角色
export function getroles(params) {
  return request({
    url: '/getroles',
    method: 'get',
    params
  })
}

// 创建用户
export function createuserapi(data) {
  return request({
    url: '/superadmin/users/creatmanager',
    method: 'post',
    data
  })
}

// 编辑用户
export function edituserstatus(data) {
  return request({
    url: '/superadmin/users/edituser',
    method: 'post',
    data
  })
}

// 获取视频列表
export function getvideolist(params) {
  return request({
    url: '/video',
    method: 'get',
    params
  })
}

// 删除视频
export function deletevideo(data) {
  return request({
    url: '/superadmin/video/delete',
    method: 'post',
    data
  })
}

// 编辑视频
export function editevideo(data) {
  return request({
    url: '/superadmin/video/delete',
    method: 'post',
    data
  })
}

// 获取分类
export function getcategory(params) {
  return request({
    url: '/video/category',
    method: 'get',
    params
  })
}

// 创建视频
export function createVideo(data) {
  return request({
    url: '/superadmin/video/create',
    method: 'post',
    data
  })
}
// 创建剧集
export function createEpisodes(data) {
  return request({
    url: '/superadmin/video/createinfo',
    method: 'post',
    data
  })
}
// 创建分类
export function createCategories(data) {
  return request({
    url: '/superadmin/video/category/create',
    method: 'post',
    data
  })
}

// 删除分类
export function deleteCategories(data) {
  return request({
    url: '/superadmin/video/category/delete',
    method: 'post',
    data
  })
}

// 获取产品
export function getProducts() {
  return request({
    url: '/superadmin/order/product',
    method: 'get'
  })
}

// 创建/编辑商品
export function createProducts(data) {
  return request({
    url: '/superadmin/order/product',
    method: 'post',
    data
  })
}

// 创建/编辑商品
export function deleteProducts(data) {
  return request({
    url: '/superadmin/order/delproduct',
    method: 'post',
    data
  })
}

// 获取订单
export function getorders(params) {
  return request({
    url: '/orders',
    method: 'get',
    params
  })
}


// 获取订单
export function getincomes(params) {
  return request({
    url: '/incomes',
    method: 'get',
    params
  })
}
