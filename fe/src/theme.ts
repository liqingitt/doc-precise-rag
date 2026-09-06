import { theme, type ThemeConfig } from 'antd'

export const appTheme: ThemeConfig = {
  algorithm: theme.defaultAlgorithm,
  token: {
    colorPrimary: '#3b6cff',
    colorInfo: '#3b6cff',
    colorLink: '#3b6cff',
    colorSuccess: '#16a34a',
    colorWarning: '#d97706',
    colorError: '#e11d48',
    colorBgLayout: '#e8eef6',
    colorBgContainer: '#ffffff',
    colorBgElevated: '#ffffff',
    colorBorder: '#d9e2ec',
    colorBorderSecondary: '#e6edf5',
    colorText: '#1e293b',
    colorTextSecondary: '#64748b',
    colorTextTertiary: '#94a3b8',
    borderRadius: 10,
    fontFamily:
      '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans SC", sans-serif',
  },
  components: {
    Layout: {
      headerBg: '#f3f6fb',
      siderBg: '#f3f6fb',
      bodyBg: '#e8eef6',
      headerHeight: 60,
      headerPadding: '0 28px',
    },
    Menu: {
      itemBg: 'transparent',
      itemSelectedBg: '#e4ecff',
      itemSelectedColor: '#2f5ef0',
      itemHoverBg: '#eaf0f8',
      itemActiveBg: '#e4ecff',
      itemColor: '#475569',
    },
    Button: {
      primaryShadow: 'none',
    },
    Table: {
      headerBg: '#f5f8fc',
      headerColor: '#475569',
      rowHoverBg: '#f7faff',
    },
    Modal: {
      headerBg: '#ffffff',
    },
  },
}
