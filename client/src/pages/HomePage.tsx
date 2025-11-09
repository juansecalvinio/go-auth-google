export const HomePage = () => {
  const handleLogin = () => {
    const redirectURL = encodeURIComponent("http://localhost:5173/profile");
    // El frontend envía el parámetro 'state' en la URL del backend
    window.location.href = `http://localhost:8080/auth/google?provider=google&state=${redirectURL}`;
  };

  return (
    <div className="relative h-screen">
      <div className="absolute inset-0">
        <div className="relative h-full w-full bg-slate-950 [&>div]:absolute [&>div]:inset-0 [&>div]:bg-[radial-gradient(circle_500px_at_50%_200px,#3e3e3e,transparent)]">
          <div></div>
        </div>
      </div>

      <div className="relative z-10 flex h-full flex-col items-center justify-center px-4">
        <div className="max-w-3xl text-center">
          <h1 className="mb-8 text-4xl font-bold tracking-tight sm:text-6xl lg:text-7xl text-white">
            Golang & React
          </h1>
          <p className="mx-auto mb-8 max-w-2xl text-lg text-slate-300">
            A modern authentication system built with a Go backend using Gin web
            framework and Goth for OAuth integration, following hexagonal
            architecture principles. And frontend crafted with React.
          </p>
          <div className="flex flex-wrap justify-center gap-4">
            <button
              onClick={handleLogin}
              className="rounded-lg px-6 py-3 font-medium  bg-rose-700 text-white hover:bg-rose-500"
            >
              Sign in with Google
            </button>
            <button className="rounded-lg border px-6 py-3 font-medium border-slate-700 bg-slate-800 text-white hover:bg-slate-700">
              <a
                className="link"
                href="https://github.com/juansecalvinio/go-auth-google"
                target="_blank"
              >
                View repository
              </a>
            </button>
          </div>
        </div>
      </div>

      {/* Footer fixed to bottom */}
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
