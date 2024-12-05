// import { useAuth } from "@/providers/auth-provider";
// import { useEffect } from "react";
// import { redirect } from "react-router-dom";

export const AuthLayout = ({ children }: { children: React.ReactNode }) => {

  return (
    <main className="flex flex-col items-center justify-center w-screen h-screen">
      <div className="flex flex-col w-full max-w-sm">
        <h1 className="text-4xl text-gray-900 font-black mb-12">Chad-RSS</h1>
        <div className="flex flex-col flex-auto">
          {children}
        </div>
      </div>
    </main>
  );
};
