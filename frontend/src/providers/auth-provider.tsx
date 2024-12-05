import { useContext, createContext, useState, useMemo, useEffect } from "react";
import { jwtDecode, JwtPayload } from 'jwt-decode';
import { redirect } from "react-router-dom";

type AuthContextValue = {
  token: string | null;
  signin: (data: any) => void;
  signup: (data: any) => void;
  signout: () => void;
  user: Token | null;
};

const AuthContext = createContext<AuthContextValue>(null!);

export const useAuth = () => {
  return useContext(AuthContext);
};

export const AuthProvider = ({ children }: React.PropsWithChildren) => {
  const [token, setToken] = useState<string | null>(localStorage.getItem("site"));

  const signin = async (data: { username: string; password: string }) => {
    try {
      const response = await fetch("http://localhost:8080/auth/signin", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });
      const res: string = await response.json();
      if (res != null) {
        setToken(res);
        return;
      }
      throw new Error("Invalid response");
    } catch (err) {
      console.error(err);
      setToken(null);
    }
  };

  const signup = async (data: { username: string; password: string }) => {
    try {
      const response = await fetch("http://localhost:8080/auth/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });
      const res: string = await response.json();
      if (res != null) {
        setToken(res);
        return;
      }
      throw new Error("Invalid response");
    } catch (err) {
      console.error(err);
      setToken(null);
    }
  };

  const signout = () => {
    setToken(null);
  };

  const user = useMemo(() => {
    if (token == null) return null;

    const decoded = jwtDecode(token) as Token;
    if (decoded == null) return null;

    return {
      id: decoded.id,
      username: decoded.username,
    };
  }, [token]);

  useEffect(() => {
    if (token != null) {
      redirect("/");
    }
  }, []);

  useEffect(() => {
    if (token == null) {
      localStorage.removeItem("site");
    } else {
      localStorage.setItem("site", token);
    }
  }, [token]);

  return (
    <AuthContext.Provider value={{ token, signin, signup, signout, user }}>
      {children}
    </AuthContext.Provider>
  );
};

type Token = JwtPayload & {
  id: string;
  username: string;
};
