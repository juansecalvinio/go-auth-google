import { Route, BrowserRouter as Router, Routes } from "react-router-dom";
import { ProfilePage } from "./pages/ProfilePage";
import { HomePage } from "./pages/HomePage";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/profile" element={<ProfilePage />} />
      </Routes>
    </Router>
  );
}

export default App;
