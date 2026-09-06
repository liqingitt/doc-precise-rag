export interface UploadURLData {
  upload_url: string
  object_key: string
  file_url: string
  filename: string
}

interface UploadURLResponse {
  code: number
  message?: string
  data?: UploadURLData
}

export async function createUploadURL(filename: string): Promise<UploadURLData> {
  const res = await fetch('/api/common/upload-url', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ filename }),
  })
  const json = (await res.json()) as UploadURLResponse
  if (!res.ok || json.code !== 0 || !json.data) {
    throw new Error(json.message || '获取上传地址失败')
  }
  return json.data
}

export function putFileToCOS(
  url: string,
  file: File,
  onProgress?: (percent: number) => void,
): Promise<void> {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest()
    xhr.open('PUT', url)
    xhr.setRequestHeader('Content-Type', 'application/pdf')
    xhr.upload.onprogress = (event) => {
      if (event.lengthComputable) {
        onProgress?.(Math.round((event.loaded / event.total) * 100))
      }
    }
    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        resolve()
        return
      }
      reject(new Error(`上传失败（${xhr.status}）`))
    }
    xhr.onerror = () => reject(new Error('上传失败，请检查 COS 跨域配置'))
    xhr.send(file)
  })
}
