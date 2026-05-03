import request from '@/utils/request'

export function userLogin(data) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}

export function userRegister(data) {
  return request({
    url: '/user/register',
    method: 'post',
    data
  })
}

export function getUserProfile() {
  return request({
    url: '/user/profile',
    method: 'get'
  })
}

export function updateUserProfile(data) {
  return request({
    url: '/user/profile',
    method: 'put',
    data
  })
}

export function changePassword(data) {
  return request({
    url: '/user/password',
    method: 'put',
    data
  })
}

export function getRoutes(params) {
  return request({
    url: '/user/routes',
    method: 'get',
    params
  })
}

export function getRouteDetail(id) {
  return request({
    url: `/user/routes/${id}`,
    method: 'get'
  })
}

export function getTickets(params) {
  return request({
    url: '/user/tickets',
    method: 'get',
    params
  })
}

export function getTicketDetail(id) {
  return request({
    url: `/user/tickets/${id}`,
    method: 'get'
  })
}

export function bookTicket(data) {
  return request({
    url: '/user/tickets/book',
    method: 'post',
    data
  })
}

export function payBooking(data) {
  return request({
    url: '/user/tickets/pay',
    method: 'post',
    data
  })
}

export function getBookings(params) {
  return request({
    url: '/user/bookings',
    method: 'get',
    params
  })
}

export function cancelBooking(id) {
  return request({
    url: `/user/bookings/${id}/cancel`,
    method: 'post'
  })
}

export function getNews(params) {
  return request({
    url: '/public/news',
    method: 'get',
    params
  })
}

export function getNewsDetail(id) {
  return request({
    url: `/public/news/${id}`,
    method: 'get'
  })
}

export function getAnnouncements(params) {
  return request({
    url: '/public/announcements',
    method: 'get',
    params
  })
}

export function getCarousels() {
  return request({
    url: '/public/carousels',
    method: 'get'
  })
}

export function getVehicleMonitors(params) {
  return request({
    url: '/user/vehicle-monitor',
    method: 'get',
    params
  })
}

export function getVehicleMonitorByPlate(plateNumber) {
  return request({
    url: `/user/vehicle-monitor/plate/${plateNumber}`,
    method: 'get'
  })
}

export function getRoadConditions(params) {
  return request({
    url: '/public/road-conditions',
    method: 'get',
    params
  })
}

export function getRoadConditionDetail(id) {
  return request({
    url: `/user/road-conditions/${id}`,
    method: 'get'
  })
}

export function getViolations(params) {
  return request({
    url: '/user/violations',
    method: 'get',
    params
  })
}

export function getViolationDetail(id) {
  return request({
    url: `/user/violations/${id}`,
    method: 'get'
  })
}

export function payViolation(id) {
  return request({
    url: `/user/violations/${id}/pay`,
    method: 'post'
  })
}

export function getFeedbacks(params) {
  return request({
    url: '/user/feedbacks',
    method: 'get',
    params
  })
}

export function getFeedbackDetail(id) {
  return request({
    url: `/user/feedbacks/${id}`,
    method: 'get'
  })
}

export function createFeedback(data) {
  return request({
    url: '/user/feedbacks',
    method: 'post',
    data
  })
}
