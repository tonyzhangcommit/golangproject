import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

export function getInfo(data) {
  return request({
    url: '/usersinfobyjwt',
    method: 'post',
    data
  })
}

export function logout() {
  return request({
    url: 'common/logout/',
    method: 'post'
  })
}
