import { FileTextOutlined } from '@ant-design/icons'
import { Layout, Menu } from 'antd'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import styles from './AdminLayout.module.less'

const pageTitle: Record<string, string> = {
  '/docs': '文档管理',
}

export default function AdminLayout() {
  const location = useLocation()
  const navigate = useNavigate()

  return (
    <Layout className={styles.layout}>
      <Layout.Sider className={styles.sider} theme="light" width={228}>
        <div className={styles.logo}>
          <span className={styles.logoMark} />
          Doc Precise RAG
        </div>
        <Menu
          mode="inline"
          selectedKeys={[location.pathname]}
          onClick={({ key }) => navigate(key)}
          items={[
            {
              key: '/docs',
              icon: <FileTextOutlined />,
              label: '文档管理',
            },
          ]}
        />
      </Layout.Sider>
      <Layout>
        <Layout.Header className={styles.header}>
          <h1 className={styles.title}>{pageTitle[location.pathname] ?? ''}</h1>
        </Layout.Header>
        <Layout.Content className={styles.content}>
          <Outlet />
        </Layout.Content>
      </Layout>
    </Layout>
  )
}
