import "~/styles/globals.css";

import { GeistSans } from "geist/font/sans";
import { type Metadata } from "next";
import { SidebarProvider, SidebarTrigger } from "~/components/ui/sidebar";
import { AppSidebar } from "~/components/app-sidebar";
import Providers from "./providers";
import { cookies } from "next/headers";

export const metadata: Metadata = {
  title: "ChadRSS",
  description: "A chad RSS reader",
  icons: [{ rel: "icon", url: "/favicon.ico" }],
};

export async function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  const cookieStore = await cookies()
  const defaultOpen = cookieStore.get("sidebar:state")?.value === "true"

  return (
    <html lang="en" className={`${GeistSans.variable}`}>
      <body>
        <Providers>
          <SidebarProvider defaultOpen={defaultOpen}>
            <AppSidebar />
            <main className="relative flex flex-col w-full h-full items-center">
              <SidebarTrigger className="absolute top-0 left-0" />
              {children}
            </main>
          </SidebarProvider>
        </Providers>
      </body>
    </html>
  );
}

export default RootLayout;
