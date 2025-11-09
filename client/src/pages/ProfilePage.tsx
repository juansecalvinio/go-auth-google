import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { ProfileCard } from "../components/ProfileCard";

interface User {
  id: string;
  name: string;
  email: string;
  picture: string;
}

export const ProfilePage = () => {
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);
  const [error, setError] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const response = await fetch("http://localhost:8080/user", {
          credentials: "include",
        });

        if (!response.ok) {
          throw new Error("Failed to fetch user data");
        }

        const userData: User = await response.json();
        setUser(userData);
      } catch (error) {
        setError(
          error instanceof Error ? error.message : "Something went wrong"
        );
      } finally {
        setLoading(false);
      }
    };

    fetchUserData();
  }, []);

  const handleLogout = async () => {
    try {
      const response = await fetch("http://localhost:8080/logout/google", {
        credentials: "include",
      });

      if (!response.ok) {
        throw new Error("Failed to logout");
      }

      // Redireccionar al home después del logout exitoso
      navigate("/");
    } catch (error) {
      setError(error instanceof Error ? error.message : "Logout failed");
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (!user) {
    return <div>No user data available</div>;
  }

  return (
    <div className="relative h-screen">
      <div className="absolute inset-0">
        <div className="relative h-full w-full bg-slate-950 [&>div]:absolute [&>div]:inset-0 [&>div]:bg-[radial-gradient(circle_500px_at_50%_200px,#3e3e3e,transparent)]">
          <div></div>
        </div>
      </div>

      <div className="relative flex place-content-center p-10 md:p-50">
        <ProfileCard
          onClick={handleLogout}
          key={user.id}
          name={user.name}
          email={user.email}
          picture={user.picture}
        />
      </div>

      {/* <div className="relative z-10 flex h-full flex-col items-center justify-center px-4">
        <div className="card bg-base-100 max-w-96 shadow-sm">
          <figure className="px-10 pt-10">
            <div className="avatar">
              <div className="rounded-full w-24">
                <img src={user.picture} alt={user.name} />
              </div>
            </div>
          </figure>

          <div className="card-body">
            <div className="flex flex-col justify-between">
              <h2 className="text-2xl font-bold">{user.name}</h2>
              <span className="text-lg">{user.email}</span>
            </div>
            <div className="mt-6">
              <button
                className="btn btn-secondary btn-block"
                onClick={handleLogout}
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </div> */}

      {/* <div className="max-w-md">
        <div className="card bg-base-100 max-w-96 shadow-sm">
          <figure className="px-10 pt-10">
            <div className="avatar">
              <div className="rounded-full w-24">
                <img src={user.picture} alt={user.name} />
              </div>
            </div>
          </figure>

          <div className="card-body">
            <div className="flex flex-col justify-between">
              <h2 className="text-2xl font-bold">{user.name}</h2>
              <span className="text-lg">{user.email}</span>
            </div>
            <div className="mt-6">
              <button
                className="btn btn-secondary btn-block"
                onClick={handleLogout}
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </div> */}

      <footer className="fixed bottom-0 left-0 w-full h-12 bg-slate-900/80 backdrop-blur-sm text-slate-400 flex items-center justify-center text-sm z-30">
        <div className="px-4">
          Created by{" "}
          <a
            className="link"
            href="https://juansecalvinio.com"
            target="_blank"
            rel="noopener noreferrer"
          >
            @juansecalvinio
          </a>
        </div>
      </footer>
    </div>
  );
};
