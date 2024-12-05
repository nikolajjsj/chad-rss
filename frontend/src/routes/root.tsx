import Layout from '@/layouts/layout';
import { useAuth } from '@/providers/auth-provider';
import { useEffect } from 'react';
import { Outlet, redirect } from 'react-router-dom';

export const Root = () => {
  const { token } = useAuth();

  useEffect(() => {
    if (token == null) {
      redirect("/signin");
    }
  }, [token]);

  return (
    <Layout>
      <div id="detail">
        <Outlet />
      </div>
    </Layout>
  );
};
