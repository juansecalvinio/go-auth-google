import reactLogo from "./assets/react.svg";
import golangLogo from "./assets/go-light-blue.svg";
import googleLogo from "./assets/google.svg";
import "./App.css";

function App() {
  const handleLogin = () => {
    const redirectURL = encodeURIComponent("http://localhost:5173/profile");
    // El frontend envía el parámetro 'state' en la URL del backend
    window.location.href = `http://localhost:8080/auth/google?state=${redirectURL}`;
  };

  return (
    <>
      <div>
        <a href="https://go.dev" target="_blank">
          <img src={golangLogo} className="logo" alt="Golang logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h2>Go + React</h2>
      <p className="read-the-docs">Login with Google Authentication</p>

      <div className="card">
        <button onClick={handleLogin} className="google-btn">
          <img src={googleLogo} className="google-icon" alt="Google logo" />
          <span>Sign in with Google</span>
        </button>
      </div>
    </>
  );
}

export default App;
