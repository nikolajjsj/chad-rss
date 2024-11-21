import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/app-sidebar"

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <SidebarProvider>
      <AppSidebar />
      <main className="relative w-full h-full">
        <div className="absolute top-0 left-0 right-0 bottom-0 w-full h-full">
          {children}
        </div>

        <SidebarTrigger className="absolute top-0 left-0" />
      </main>
    </SidebarProvider>
  )
}
