import { useState } from "react";
import { Card, CardContent } from "../Components/Card";
import { useNavigate } from "react-router-dom";
import Button from "../Components/CommonButton";
import { useAuthStore } from "../hooks";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();
  const { login } = useAuthStore();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);
    try {
      await login({ username: email, password });
      navigate("/create-workspace");
    } catch (err: any) {
      setError(
        err.response?.data?.message ||
          "Error al iniciar sesión. Inténtalo de nuevo."
      );
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen w-full items-center justify-center">
      <div className="w-full max-w-sm">
        <Card className="p-4">
          <div className="select-none">
            <h1 className="text-2xl ">Inicio de sesión</h1>
            <span className="text-gray-400 text-sm">
              Coloque su nombre de usuario y su contraseña para continuar
            </span>
          </div>
          <CardContent>
            <form onSubmit={handleLogin}>
              <div className="flex flex-col gap-6">
                <div className="grid gap-2">
                  <label htmlFor="userName" className="text-gray-500 text-sm">
                    Nombre de usuario:
                  </label>
                  <input
                    id="userName"
                    type="text"
                    className={`flex-1 px-2 py-1 border border-gray-400 rounded text-sm focus:outline-none focus:ring-2 focus:ring-cyan-500`}
                    required
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                </div>
                <div className="grid gap-2">
                  <label htmlFor="password" className="text-gray-500 text-sm">
                    Contraseña:
                  </label>
                  <input
                    id="password"
                    type="password"
                    className={`flex-1 px-2 py-1 border border-gray-400 rounded text-sm focus:outline-none focus:ring-2 focus:ring-cyan-500`}
                    required
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                  />
                </div>
                {error && <p className="text-sm text-destructive">{error}</p>}
                <Button
                  type="submit"
                  className="w-full bg-emerald-500 hover:bg-emerald-600 text-white rounded-lg"
                  disabled={isLoading}
                >
                  {isLoading ? "Iniciando sesión..." : "Iniciar sesión"}
                </Button>
              </div>
              <div className="mt-4 text-center text-sm">
                No contas con una cuenta?
                <Button
                  page="auth/sign-up"
                  className="text-black! rounded-lg border-none underline"
                >
                  Regístrate
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
