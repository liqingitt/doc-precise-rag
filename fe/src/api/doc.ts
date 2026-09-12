import { postJSON } from './http'

export interface DocOriginFileItem {
  id: string
  object_key: string
  doc_title: string
  create_time: string
  update_time: string
}

export interface DocOriginFilePage {
  total: number
  list: DocOriginFileItem[] | null
}

export function queryDocOriginFileList(page: number, pageSize: number) {
  return postJSON<DocOriginFilePage>('/api/doc/query-doc-origin-file-list', {
    page,
    page_size: pageSize,
  })
}

export function addDocOriginFile(objectKey: string, docTitle: string) {
  return postJSON<null>('/api/doc/add-doc-origin-file', {
    object_key: objectKey,
    doc_title: docTitle,
  })
}

export function analysisDocOriginFile(id: string) {
  return postJSON<null>('/api/doc/analysis-doc-origin-file', {
    id,
  })
}

export function deleteDocOriginFile(id: string) {
  return postJSON<null>('/api/doc/delete-doc-origin-file', {
    id,
  })
}
