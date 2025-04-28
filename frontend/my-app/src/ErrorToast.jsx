// ErrorToast.jsx
import { Snackbar, Alert } from '@mui/material';

export default function ErrorToast({ error, onClose }) {
  return (
    <Snackbar open={!!error} autoHideDuration={6000} onClose={onClose}>
      <Alert severity="error">{error}</Alert>
    </Snackbar>
  );
}