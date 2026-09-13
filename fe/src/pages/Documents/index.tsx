import { InboxOutlined, PlusOutlined, ReloadOutlined } from '@ant-design/icons'
import { App, Button, Input, Modal, Popconfirm, Space, Spin, Table, Typography, Upload } from 'antd'
import type { TableColumnsType } from 'antd'
import { useCallback, useEffect, useState } from 'react'
import Markdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { createUploadURL, putFileToCOS } from '../../api/upload'
import {
  addDocOriginFile,
  analysisDocOriginFile,
  deleteDocOriginFile,
  getTempFileUrlByObjectKey,
  queryDocOriginFileList,
  type DocOriginFileItem,
} from '../../api/doc'
import styles from './index.module.less'

function isPDF(file: File) {
  return file.type === 'application/pdf' || file.name.toLowerCase().endsWith('.pdf')
}

function titleFromFilename(name: string) {
  return name.replace(/\.pdf$/i, '')
}

function formatTime(value?: string) {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleString('zh-CN', { hour12: false })
}

function EllipsisCell({ text, copyable }: { text?: string; copyable?: boolean }) {
  const value = text || '-'
  return (
    <Typography.Text
      className={styles.cellText}
      ellipsis={{ tooltip: value }}
      copyable={copyable ? { text: value } : undefined}
    >
      {value}
    </Typography.Text>
  )
}

export default function DocumentsPage() {
  const { message } = App.useApp()
  const [list, setList] = useState<DocOriginFileItem[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [loading, setLoading] = useState(false)
  const [analyzingKeys, setAnalyzingKeys] = useState<string[]>([])
  const [deletingKeys, setDeletingKeys] = useState<string[]>([])

  const [modalOpen, setModalOpen] = useState(false)
  const [file, setFile] = useState<File | null>(null)
  const [docTitle, setDocTitle] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [uploadPercent, setUploadPercent] = useState<number | null>(null)

  const [previewOpen, setPreviewOpen] = useState(false)
  const [previewTitle, setPreviewTitle] = useState('')
  const [previewMarkdown, setPreviewMarkdown] = useState('')
  const [previewLoading, setPreviewLoading] = useState(false)

  const loadList = useCallback(
    async (nextPage = page, nextPageSize = pageSize) => {
      setLoading(true)
      try {
        const data = await queryDocOriginFileList(nextPage, nextPageSize)
        setList(data.list ?? [])
        setTotal(data.total ?? 0)
      } catch (error) {
        message.error(error instanceof Error ? error.message : '查询文档失败')
      } finally {
        setLoading(false)
      }
    },
    [message, page, pageSize],
  )

  useEffect(() => {
    void loadList()
  }, [loadList])

  const resetModal = () => {
    setFile(null)
    setDocTitle('')
    setUploadPercent(null)
    setSubmitting(false)
  }

  const handleAdd = async () => {
    const title = docTitle.trim()
    if (!file) {
      message.warning('请先选择 PDF 文件')
      return
    }
    if (!title) {
      message.warning('请填写文档标题')
      return
    }

    setSubmitting(true)
    setUploadPercent(0)
    try {
      const signed = await createUploadURL(file.name)
      await putFileToCOS(signed.upload_url, file, (percent) => {
        setUploadPercent(percent)
      })
      await addDocOriginFile(signed.object_key, title)
      message.success('新增文档成功')
      setModalOpen(false)
      resetModal()
      if (page === 1) {
        await loadList(1, pageSize)
      } else {
        setPage(1)
      }
    } catch (error) {
      message.error(error instanceof Error ? error.message : '新增文档失败')
    } finally {
      setSubmitting(false)
      setUploadPercent(null)
    }
  }

  const handleAnalysis = async (record: DocOriginFileItem) => {
    setAnalyzingKeys((prev) => [...prev, record.id])
    try {
      await analysisDocOriginFile(record.id)
      message.success(`已解析「${record.doc_title}」`)
      await loadList(page, pageSize)
    } catch (error) {
      message.error(error instanceof Error ? error.message : '解析失败')
    } finally {
      setAnalyzingKeys((prev) => prev.filter((key) => key !== record.id))
    }
  }

  const handleDelete = async (record: DocOriginFileItem) => {
    setDeletingKeys((prev) => [...prev, record.id])
    try {
      await deleteDocOriginFile(record.id)
      message.success(`已删除「${record.doc_title}」`)
      const shouldGoPrevPage = list.length === 1 && page > 1
      const nextPage = shouldGoPrevPage ? page - 1 : page
      if (shouldGoPrevPage) {
        setPage(nextPage)
      }
      await loadList(nextPage, pageSize)
    } catch (error) {
      message.error(error instanceof Error ? error.message : '删除失败')
    } finally {
      setDeletingKeys((prev) => prev.filter((key) => key !== record.id))
    }
  }

  const handleViewParsed = async (record: DocOriginFileItem) => {
    const objectKey = record.analyze_doc_object_key?.trim()
    if (!objectKey) return

    setPreviewTitle(record.doc_title)
    setPreviewMarkdown('')
    setPreviewOpen(true)
    setPreviewLoading(true)
    try {
      const { url } = await getTempFileUrlByObjectKey(objectKey)
      const res = await fetch(url)
      if (!res.ok) {
        throw new Error(`拉取解析文档失败（${res.status}）`)
      }
      const markdown = await res.text()
      setPreviewMarkdown(markdown)
    } catch (error) {
      message.error(error instanceof Error ? error.message : '查看解析文档失败')
      setPreviewOpen(false)
    } finally {
      setPreviewLoading(false)
    }
  }

  const columns: TableColumnsType<DocOriginFileItem> = [
    {
      title: '文档标题',
      dataIndex: 'doc_title',
      ellipsis: { showTitle: false },
      render: (value: string) => <EllipsisCell text={value} />,
    },
    {
      title: '对象键',
      dataIndex: 'object_key',
      ellipsis: { showTitle: false },
      render: (value: string) => <EllipsisCell text={value} copyable />,
    },
    {
      title: '创建时间',
      dataIndex: 'create_time',
      width: 200,
      render: formatTime,
    },
    {
      title: '操作',
      key: 'action',
      width: 280,
      render: (_, record) => {
        const hasParsed = Boolean(record.analyze_doc_object_key?.trim())
        return (
          <Space size={0}>
            <Button
              type="link"
              disabled={!hasParsed}
              onClick={() => void handleViewParsed(record)}
            >
              查看已解析
            </Button>
            <Button
              type="link"
              onClick={() => void handleAnalysis(record)}
              loading={analyzingKeys.includes(record.id)}
            >
              开始解析
            </Button>
            <Popconfirm
              title="删除文档"
              description={`确认删除「${record.doc_title}」？删除后不可恢复。`}
              okText="删除"
              cancelText="取消"
              okButtonProps={{ danger: true }}
              onConfirm={() => handleDelete(record)}
            >
              <Button type="link" danger loading={deletingKeys.includes(record.id)}>
                删除
              </Button>
            </Popconfirm>
          </Space>
        )
      },
    },
  ]

  return (
    <div className={styles.page}>
      <div className={styles.toolbar}>
        <div>
          <div className="eyebrow">Knowledge Base</div>
          <h2 className={`${styles.heading} page-title`}>原始文档</h2>
          <p className={styles.desc}>上传 PDF 后直传对象存储，再登记文档并启动解析。</p>
        </div>
        <Space size={8}>
          <Button icon={<ReloadOutlined />} loading={loading} onClick={() => void loadList(page, pageSize)}>
            刷新
          </Button>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalOpen(true)}>
            新增文档
          </Button>
        </Space>
      </div>

      <Table<DocOriginFileItem>
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={list}
        tableLayout="fixed"
        bordered
        pagination={{
          current: page,
          pageSize,
          total,
          showSizeChanger: true,
          showTotal: (value) => `共 ${value} 条`,
          onChange: (nextPage, nextPageSize) => {
            setPage(nextPage)
            setPageSize(nextPageSize)
          },
        }}
      />

      <Modal
        title="新增文档"
        open={modalOpen}
        okText={uploadPercent != null ? `上传中 ${uploadPercent}%` : '确认添加'}
        confirmLoading={submitting}
        onOk={() => void handleAdd()}
        onCancel={() => {
          if (submitting) return
          setModalOpen(false)
          resetModal()
        }}
        destroyOnHidden
      >
        <Space direction="vertical" size={16} className={styles.modalBody}>
          <Upload.Dragger
            accept="application/pdf,.pdf"
            maxCount={1}
            disabled={submitting}
            fileList={
              file
                ? [
                    {
                      uid: 'selected-pdf',
                      name: file.name,
                      status: submitting ? 'uploading' : 'done',
                      percent: uploadPercent ?? undefined,
                    },
                  ]
                : []
            }
            beforeUpload={(nextFile) => {
              if (!isPDF(nextFile)) {
                message.error('仅支持 PDF 文件')
                return Upload.LIST_IGNORE
              }
              setFile(nextFile)
              setDocTitle(titleFromFilename(nextFile.name))
              return false
            }}
            onRemove={() => {
              if (submitting) return false
              setFile(null)
              setDocTitle('')
              return true
            }}
          >
            <p className="ant-upload-drag-icon">
              <InboxOutlined />
            </p>
            <p className="ant-upload-text">点击或拖拽 PDF 到这里</p>
            <p className="ant-upload-hint">文件将直传到 COS，不会经过业务服务器中转</p>
          </Upload.Dragger>
          <div>
            <div className={styles.fieldLabel}>文档标题</div>
            <Input
              value={docTitle}
              maxLength={120}
              placeholder="默认使用原始文件名（不含后缀）"
              disabled={submitting}
              onChange={(event) => setDocTitle(event.target.value)}
            />
          </div>
        </Space>
      </Modal>

      <Modal
        title={previewTitle ? `已解析文档 · ${previewTitle}` : '已解析文档'}
        open={previewOpen}
        width={880}
        footer={null}
        onCancel={() => {
          if (previewLoading) return
          setPreviewOpen(false)
          setPreviewMarkdown('')
        }}
        destroyOnHidden
      >
        <Spin spinning={previewLoading}>
          <div className={styles.markdownPreview}>
            {previewMarkdown ? (
              <Markdown remarkPlugins={[remarkGfm]}>{previewMarkdown}</Markdown>
            ) : (
              <div className={styles.markdownEmpty}>{previewLoading ? '加载中…' : '暂无内容'}</div>
            )}
          </div>
        </Spin>
      </Modal>
    </div>
  )
}
