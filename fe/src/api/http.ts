interface BaseResp<T> {
  success?: boolean
  code?: number
  message?: string
  data?: T
}

export async function postJSON<T>(url: string, body: unknown = {}): Promise<T> {
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })

  let json: BaseResp<T>
  try {
    json = (await res.json()) as BaseResp<T>
  } catch {
    throw new Error('服务响应异常')
  }

  if (!res.ok || json.success === false) {
    throw new Error(json.message || '请求失败')
  }

  return json.data as T
}
