import request from "@/utils/request"

export const dictApi = {
  // Dict Type APIs
  getDictTypeList(params: any) {
    return request({
      url: '/api/dict/type/list',
      method: 'get',
      params
    })
  },
  addDictType(data: any) {
    return request({
      url: '/api/dict/type/add',
      method: 'post',
      data
    })
  },
  updateDictType(data: any) {
    return request({
      url: '/api/dict/type/update',
      method: 'put',
      data
    })
  },
  deleteDictType(data: { ids: number[] }) {
    return request({
      url: '/api/dict/type/delete',
      method: 'delete',
      data
    })
  },

  // Dict Data APIs
  getDictDataList(params: any) {
    return request({
      url: '/api/dict/data/list',
      method: 'get',
      params
    })
  },
  addDictData(data: any) {
    return request({
      url: '/api/dict/data/add',
      method: 'post',
      data
    })
  },
  updateDictData(data: any) {
    return request({
      url: '/api/dict/data/update',
      method: 'put',
      data
    })
  },
  deleteDictData(data: { ids: number[] }) {
    return request({
      url: '/api/dict/data/delete',
      method: 'delete',
      data
    })
  }
}
