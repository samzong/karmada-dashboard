import './App.css';
import Router from './routes';
import { Helmet, HelmetProvider } from 'react-helmet-async';
import { ConfigProvider } from 'antd';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import AuthProvider from '@/components/auth';
import { getAntdLocale } from '@/utils/i18n.tsx';

const queryClient = new QueryClient();

function App() {
  return (
    <ConfigProvider
      locale={getAntdLocale()}
      theme={{
        components: {
          Layout: {
            siderBg: '#ffffff',
          },
        },
      }}
    >
      <QueryClientProvider client={queryClient}>
        <AuthProvider>
          <HelmetProvider>
            <Helmet>
              <title>Karmada Dashboard</title>
              <link
                rel="apple-touch-icon"
                sizes="180x180"
                href="/apple-touch-icon.png"
              />

              <link
                rel="icon"
                type="image/png"
                sizes="16x16"
                href="/favicon-16x16.png"
              />

              <link
                rel="icon"
                type="image/png"
                sizes="32x32"
                href="/favicon-32x32.png"
              />

              <link
                rel="shortcut icon"
                type="image/x-icon"
                href="/favicon.ico"
              />
            </Helmet>
            <Router />
          </HelmetProvider>
        </AuthProvider>
      </QueryClientProvider>
    </ConfigProvider>
  );
}

export default App;
