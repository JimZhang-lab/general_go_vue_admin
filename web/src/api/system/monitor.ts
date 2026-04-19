import request from "@/utils/request"

export const monitorApi = {
  // Get Server Monitor Info
  getServerInfo() {
    return request({
      url: '/api/monitor/server',
      method: 'get'
    })
  }
}
