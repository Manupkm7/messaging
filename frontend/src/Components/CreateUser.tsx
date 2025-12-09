import { useState } from "react";

//import { UserPlus } from "lucide-react";
import Modal from "./Modal";
import Button from "./CommonButton";

interface CreateUserModalProps {
  workspaceId: string;
  onUserCreated: () => void;
}

export function CreateUserModal({
  workspaceId,
  onUserCreated,
}: CreateUserModalProps) {
  const [open, setOpen] = useState(false);
  const [email, setEmail] = useState("");
  // const [displayName, setDisplayName] = useState("");
  // const [password, setPassword] = useState("");
  // const [isAdmin, setIsAdmin] = useState(false);
  const [role, setRole] = useState<"admin" | "member">("member");
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);
    onUserCreated();
    console.debug("User created:", workspaceId);
  };

  return (
    <Modal show={open} onClose={() => setOpen(false)}>
      <div className="sm:max-w-md">
        <div>
          <h1>Add User to Workspace</h1>
          <span>
            Add an existing user to this workspace. They must have an account
            already.
          </span>
        </div>
        <form onSubmit={handleSubmit}>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <label htmlFor="email">Email</label>
              <input
                id="email"
                type="email"
                placeholder="user@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>

            <div className="space-y-2">
              <label htmlFor="role">Workspace Role</label>
              <select
                id="role"
                value={role}
                onChange={(e) => setRole(e.target.value as "admin" | "member")}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background"
              >
                <option value="member">Member</option>
                <option value="admin">Workspace Admin</option>
              </select>
            </div>

            <div className="flex items-center justify-between space-x-2">
              <label htmlFor="isAdmin" className="flex-1">
                System Administrator
              </label>
              {/*  <Switch
                id="isAdmin"
                checked={isAdmin}
                onCheckedChange={setIsAdmin}
              />
            */}
            </div>
            <p className="text-xs text-muted-foreground">
              System admins can manage all workspaces and create users
            </p>

            {error && <p className="text-sm text-destructive">{error}</p>}
          </div>

          <div>
            <Button type="button" onClick={() => setOpen(false)}>
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading ? "Adding..." : "Add User"}
            </Button>
          </div>
        </form>
      </div>
    </Modal>
  );
}
