import { Card, CardContent } from "../Components/Card";
import Button from "../Components/CommonButton";
import type React from "react";

import { useState } from "react";
import { useNavigate } from "react-router-dom";

export default function CreateWorkspacePage() {
  const [name, setName] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  const generateSlug = (name: string) => {
    return name
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, "-")
      .replace(/(^-|-$)/g, "");
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);
    navigate("/");
  };

  return (
    <div className="flex min-h-screen w-full items-center justify-center">
      <div className="w-full max-w-md">
        <Card className="p-4">
          <div className="select-none">
            <h1 className="text-2xl">Crear Grupo de trabajo</h1>
            <span className="text-gray-400 text-sm">
              Crea un nuevo grupo de trabajo para comenzar a colaborar con tu
              equipo.
            </span>
          </div>
          <CardContent>
            <form onSubmit={handleCreate}>
              <div className="flex flex-col gap-6">
                <div className="grid gap-2">
                  <label htmlFor="name" className="text-gray-500 text-sm">Grupo de trabajo</label>
                  <input
                    id="name"
                    type="text"
                    className={`flex-1 px-2 py-1 border border-gray-400 rounded text-sm focus:outline-none focus:ring-2 focus:ring-cyan-500`}
                    required
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                  />
                  {name && (
                    <p className="text-xs text-muted-foreground">
                      Slug: {generateSlug(name)}
                    </p>
                  )}
                </div>
                {error && <p className="text-sm text-destructive">{error}</p>}
                <Button type="submit" className="w-full" disabled={isLoading}>
                  {isLoading ? "Creating..." : "Create Workspace"}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
