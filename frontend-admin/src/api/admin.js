import request from '@/utils/request'

export function adminLogin(data) {
  return request({
    url: '/admin/login',
    method: 'post',
    data
  })
}

export function getAdminProfile() {
  return request({
    url: '/admin/profile',
    method: 'get'
  })
}

export function changePassword(data) {
  return request({
    url: '/admin/password',
    method: 'put',
    data
  })
}

export function getUsers(params) {
  return request({
    url: '/admin/users',
    method: 'get',
    params
  })
}

export function updateUserStatus(id, status) {
  return request({
    url: `/admin/users/${id}/status`,
    method: 'put',
    data: { status }
  })
}

export function getAdmins(params) {
  return request({
    url: '/admin/admins',
    method: 'get',
    params
  })
}

export function createAdmin(data) {
  return request({
    url: '/admin/admins',
    method: 'post',
    data
  })
}

export function getVehicles(params) {
  return request({
    url: '/admin/vehicles',
    method: 'get',
    params
  })
}

export function createVehicle(data) {
  return request({
    url: '/admin/vehicles',
    method: 'post',
    data
  })
}

export function updateVehicle(id, data) {
  return request({
    url: `/admin/vehicles/${id}`,
    method: 'put',
    data
  })
}

export function deleteVehicle(id) {
  return request({
    url: `/admin/vehicles/${id}`,
    method: 'delete'
  })
}

export function getVehicleMonitors(params) {
  return request({
    url: '/admin/vehicle-monitors',
    method: 'get',
    params
  })
}

export function createVehicleMonitor(data) {
  return request({
    url: '/admin/vehicle-monitors',
    method: 'post',
    data
  })
}

export function updateVehicleMonitor(id, data) {
  return request({
    url: `/admin/vehicle-monitors/${id}`,
    method: 'put',
    data
  })
}

export function deleteVehicleMonitor(id) {
  return request({
    url: `/admin/vehicle-monitors/${id}`,
    method: 'delete'
  })
}

export function getRoadConditions(params) {
  return request({
    url: '/admin/road-conditions',
    method: 'get',
    params
  })
}

export function getRoadConditionDetail(id) {
  return request({
    url: `/admin/road-conditions/${id}`,
    method: 'get'
  })
}

export function createRoadCondition(data) {
  return request({
    url: '/admin/road-conditions',
    method: 'post',
    data
  })
}

export function updateRoadCondition(id, data) {
  return request({
    url: `/admin/road-conditions/${id}`,
    method: 'put',
    data
  })
}

export function deleteRoadCondition(id) {
  return request({
    url: `/admin/road-conditions/${id}`,
    method: 'delete'
  })
}

export function getViolations(params) {
  return request({
    url: '/admin/violations',
    method: 'get',
    params
  })
}

export function getViolationDetail(id) {
  return request({
    url: `/admin/violations/${id}`,
    method: 'get'
  })
}

export function createViolation(data) {
  return request({
    url: '/admin/violations',
    method: 'post',
    data
  })
}

export function updateViolation(id, data) {
  return request({
    url: `/admin/violations/${id}`,
    method: 'put',
    data
  })
}

export function deleteViolation(id) {
  return request({
    url: `/admin/violations/${id}`,
    method: 'delete'
  })
}

export function getTickets(params) {
  return request({
    url: '/admin/tickets',
    method: 'get',
    params
  })
}

export function getTicketDetail(id) {
  return request({
    url: `/admin/tickets/${id}`,
    method: 'get'
  })
}

export function createTicket(data) {
  return request({
    url: '/admin/tickets',
    method: 'post',
    data
  })
}

export function updateTicket(id, data) {
  return request({
    url: `/admin/tickets/${id}`,
    method: 'put',
    data
  })
}

export function deleteTicket(id) {
  return request({
    url: `/admin/tickets/${id}`,
    method: 'delete'
  })
}

export function getBookings(params) {
  return request({
    url: '/admin/bookings',
    method: 'get',
    params
  })
}

export function getBookingDetail(id) {
  return request({
    url: `/admin/bookings/${id}`,
    method: 'get'
  })
}

export function cancelBooking(id) {
  return request({
    url: `/admin/bookings/${id}/cancel`,
    method: 'post'
  })
}

export function getFeedbacks(params) {
  return request({
    url: '/admin/feedbacks',
    method: 'get',
    params
  })
}

export function getFeedbackDetail(id) {
  return request({
    url: `/admin/feedbacks/${id}`,
    method: 'get'
  })
}

export function replyFeedback(id, reply) {
  return request({
    url: `/admin/feedbacks/${id}/reply`,
    method: 'put',
    data: { reply }
  })
}

export function updateFeedback(id, data) {
  return request({
    url: `/admin/feedbacks/${id}`,
    method: 'put',
    data
  })
}

export function deleteFeedback(id) {
  return request({
    url: `/admin/feedbacks/${id}`,
    method: 'delete'
  })
}

export function getTrafficFlows(params) {
  return request({
    url: '/admin/traffic-flows',
    method: 'get',
    params
  })
}

export function getTrafficFlowDetail(id) {
  return request({
    url: `/admin/traffic-flows/${id}`,
    method: 'get'
  })
}

export function createTrafficFlow(data) {
  return request({
    url: '/admin/traffic-flows',
    method: 'post',
    data
  })
}

export function updateTrafficFlow(id, data) {
  return request({
    url: `/admin/traffic-flows/${id}`,
    method: 'put',
    data
  })
}

export function deleteTrafficFlow(id) {
  return request({
    url: `/admin/traffic-flows/${id}`,
    method: 'delete'
  })
}

export function getSignalLights(params) {
  return request({
    url: '/admin/signal-lights',
    method: 'get',
    params
  })
}

export function getSignalLightDetail(id) {
  return request({
    url: `/admin/signal-lights/${id}`,
    method: 'get'
  })
}

export function createSignalLight(data) {
  return request({
    url: '/admin/signal-lights',
    method: 'post',
    data
  })
}

export function updateSignalLight(id, data) {
  return request({
    url: `/admin/signal-lights/${id}`,
    method: 'put',
    data
  })
}

export function controlSignalLight(id, data) {
  return request({
    url: `/admin/signal-lights/${id}/control`,
    method: 'put',
    data
  })
}

export function deleteSignalLight(id) {
  return request({
    url: `/admin/signal-lights/${id}`,
    method: 'delete'
  })
}

export function getRoutes(params) {
  return request({
    url: '/admin/routes',
    method: 'get',
    params
  })
}

export function getRouteDetail(id) {
  return request({
    url: `/admin/routes/${id}`,
    method: 'get'
  })
}

export function createRoute(data) {
  return request({
    url: '/admin/routes',
    method: 'post',
    data
  })
}

export function updateRoute(id, data) {
  return request({
    url: `/admin/routes/${id}`,
    method: 'put',
    data
  })
}

export function deleteRoute(id) {
  return request({
    url: `/admin/routes/${id}`,
    method: 'delete'
  })
}

export function getNews(params) {
  return request({
    url: '/admin/news',
    method: 'get',
    params
  })
}

export function getNewsDetail(id) {
  return request({
    url: `/admin/news/${id}`,
    method: 'get'
  })
}

export function createNews(data) {
  return request({
    url: '/admin/news',
    method: 'post',
    data
  })
}

export function updateNews(id, data) {
  return request({
    url: `/admin/news/${id}`,
    method: 'put',
    data
  })
}

export function deleteNews(id) {
  return request({
    url: `/admin/news/${id}`,
    method: 'delete'
  })
}

export function getAnnouncements(params) {
  return request({
    url: '/admin/announcements',
    method: 'get',
    params
  })
}

export function getAnnouncementDetail(id) {
  return request({
    url: `/admin/announcements/${id}`,
    method: 'get'
  })
}

export function createAnnouncement(data) {
  return request({
    url: '/admin/announcements',
    method: 'post',
    data
  })
}

export function updateAnnouncement(id, data) {
  return request({
    url: `/admin/announcements/${id}`,
    method: 'put',
    data
  })
}

export function deleteAnnouncement(id) {
  return request({
    url: `/admin/announcements/${id}`,
    method: 'delete'
  })
}

export function getCarousels(params) {
  return request({
    url: '/admin/carousels',
    method: 'get',
    params
  })
}

export function createCarousel(data) {
  return request({
    url: '/admin/carousels',
    method: 'post',
    data
  })
}

export function updateCarousel(id, data) {
  return request({
    url: `/admin/carousels/${id}`,
    method: 'put',
    data
  })
}

export function deleteCarousel(id) {
  return request({
    url: `/admin/carousels/${id}`,
    method: 'delete'
  })
}

export function updateAdmin(id, data) {
  return request({
    url: `/admin/admins/${id}`,
    method: 'put',
    data
  })
}

export function deleteAdmin(id) {
  return request({
    url: `/admin/admins/${id}`,
    method: 'delete'
  })
}

export function updateAdminProfile(data) {
  return request({
    url: '/admin/profile',
    method: 'put',
    data
  })
}

export function changeAdminPassword(data) {
  return request({
    url: '/admin/change-password',
    method: 'post',
    data
  })
}
