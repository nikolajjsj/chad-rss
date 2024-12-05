import { StrictMode } from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { createRoot } from 'react-dom/client'
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import './index.css'
import { AuthProvider } from './providers/auth-provider.tsx';
import { Root } from './routes/root.tsx';
import ErrorPage from './error-page.tsx';
import { Signup } from './routes/signup.tsx';
import { Signin } from './routes/signin.tsx';

// Create a client
const queryClient = new QueryClient()

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    errorElement: <ErrorPage />,
    children: [
      {
        path: "feeds/:feed",
        element: <div>Feed</div>,
      },
      {
        path: "feeds/:feed/articles/:article",
        element: <div>Feed Article</div>,
      },
      {
        path: "settings",
        element: <div>Settings</div>,
      },
    ],
  },
  {
    path: "/signup",
    element: <Signup />,
  },
  {
    path: "/signin",
    element: <Signin />,
  },
], {
  future: {
    v7_fetcherPersist: true,
    v7_relativeSplatPath: true,
    v7_partialHydration: true,
    v7_skipActionErrorRevalidation: true,
    v7_normalizeFormMethod: true,
  },
});

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <RouterProvider
          router={router}
          future={{
            v7_startTransition: true,
          }}
        />
      </AuthProvider>
    </QueryClientProvider>
  </StrictMode>,
)

