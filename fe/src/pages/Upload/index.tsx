import { InboxOutlined } from '@ant-design/icons'
import { Alert, App, Typography, Upload } from 'antd'
import { useState } from 'react'
import type { UploadFile, UploadProps } from 'antd'
import { createUploadURL, putFileToCOS, type UploadURLData } from '../../api/upload'
import styles from './index.module.less'

function isPDF(file: File) {
  return file.type === 'application/pdf' || file.name.toLowerCase().endsWith('.pdf')
}

export default function UploadPage() {
  const { message } = App.useApp()
  const [fileList, setFileList] = useState<UploadFile<UploadURLData>[]>([])
  const [result, setResult] = useState<UploadURLData | null>(null)

  const customRequest: UploadProps<UploadURLData>['customRequest'] = async (options) => {
    const { file, onError, onProgress, onSuccess } = options
    const raw = file as File
    try {
      const signed = await createUploadURL(raw.name)
      await putFileToCOS(signed.upload_url, raw, (percent) => {
        onProgress?.({ percent })
      })
      setResult(signed)
      onSuccess?.(signed)
      message.success('上传成功')
    } catch (error) {
      const err = error instanceof Error ? error : new Error('上传失败')
      onError?.(err)
      message.error(err.message)
    }
  }

  return (
    <div className={styles.page}>
      <p className={styles.desc}>
        仅支持 PDF。文件将直传到 COS 目录 <Typography.Text code>doc-anchor/pdf</Typography.Text>
        ，对象名由后端用 UUID 做前缀保证唯一。
      </p>
      <Upload.Dragger
        accept="application/pdf,.pdf"
        maxCount={1}
        fileList={fileList}
        onChange={({ fileList: next }) => setFileList(next)}
        beforeUpload={(file) => {
          if (!isPDF(file)) {
            message.error('仅支持 PDF 文件')
            return Upload.LIST_IGNORE
          }
          setResult(null)
          return true
        }}
        customRequest={customRequest}
      >
        <p className="ant-upload-drag-icon">
          <InboxOutlined />
        </p>
        <p className="ant-upload-text">点击或拖拽 PDF 到这里上传</p>
        <p className="ant-upload-hint">前端直传，不会经过业务服务器中转文件内容</p>
      </Upload.Dragger>
      {result ? (
        <Alert
          className={styles.result}
          type="success"
          showIcon
          title="上传完成"
          description={
            <div>
              <div>对象键：{result.object_key}</div>
              <div>
                访问地址：
                <Typography.Link className={styles.link} href={result.file_url} target="_blank">
                  {result.file_url}
                </Typography.Link>
              </div>
            </div>
          }
        />
      ) : null}
    </div>
  )
}
