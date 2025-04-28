import { useEffect } from 'react';
import axios from 'axios';

export default function AuthScreen() {
  useEffect(() => {
    const script = document.createElement('script');
    script.src = 'https://telegram.org/js/telegram-widget.js?22';
    script.async = true;
    script.dataset.telegramLogin = "YourBotName";
    script.dataset.size = "large";
    script.dataset.authUrl = "http://localhost:8080/api/auth";
    document.body.appendChild(script);

    return () => script.remove();
  }, []);

  return <div style={{ marginTop: '20px' }} />;
}