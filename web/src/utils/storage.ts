/*
 * @Author: JimZhang
 * @Date: 2025-07-26 20:55:08
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-26 21:01:16
 * @FilePath: /go-vue-general-admin/web/src/utils/storaage.ts
 * @Description: storaage 封装，保存后台传过来的信息
 * 
 */
interface StorageData {
  [key: string]: any;
}

export default {
  getStorage(): StorageData {
    return JSON.parse(window.localStorage.getItem(import.meta.env.VITE_NAME_SPACE as string) || "{}");
  },

  setItem(key: string, val: any): void {
    const storage = this.getStorage();
    storage[key] = val;
    window.localStorage.setItem(import.meta.env.VITE_NAME_SPACE as string, JSON.stringify(storage));
  },

  getItem(key: string): any {
    return this.getStorage()[key];
  },

  clearItem(key: string): void {
    const storage = this.getStorage();
    delete storage[key];
    window.localStorage.setItem(import.meta.env.VITE_NAME_SPACE as string, JSON.stringify(storage));
  },

  removeItem(key: string): void {
    this.clearItem(key);
  },

  clearAll(): void {
    window.localStorage.clear();
  }
};