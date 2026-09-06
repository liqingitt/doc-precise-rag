import { App as AntdApp, ConfigProvider } from 'antd'
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import zhCN from 'antd/locale/zh_CN'
import AdminLayout from './layouts/AdminLayout'
import DocumentsPage from './pages/Documents'
import { appTheme } from './theme'

export default function App() {
  return (
    <ConfigProvider locale={zhCN} theme={appTheme}>
      <AntdApp>
        <BrowserRouter>
          <Routes>
            <Route element={<AdminLayout />}>
              <Route path="/docs" element={<DocumentsPage />} />
              <Route path="/" element={<Navigate to="/docs" replace />} />
              <Route path="/upload" element={<Navigate to="/docs" replace />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </AntdApp>
    </ConfigProvider>
  )
}
