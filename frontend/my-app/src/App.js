import { BrowserRouter, Routes, Route } from 'react-router-dom';
import AuthScreen from './AuthScreen';
import MainScreen from './MainScreen';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<AuthScreen />} />
        <Route path="/main" element={<MainScreen />} />
      </Routes>
    </BrowserRouter>
  );
}   