/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /web/src/utils/authForms.ts
 * @Description: 
 * 
 */
export const USERNAME_PATTERN = /^[A-Za-z0-9_]{4,20}$/
export const PHONE_PATTERN = /^1[3-9]\d{9}$/
export const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
export const PASSWORD_PATTERN = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d@$!%*#?&_.-]{6,32}$/
export const CAPTCHA_PATTERN = /^[A-Za-z0-9]{4,6}$/

export const getRequestErrorMessage = (error: unknown, fallback: string): string => {
  const maybeError = error as {
    message?: string
    response?: {
      data?: {
        message?: string
      }
    }
  }

  return maybeError.response?.data?.message || maybeError.message || fallback
}
