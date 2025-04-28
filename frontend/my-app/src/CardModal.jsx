import { Dialog, DialogTitle, DialogContent, DialogActions, Button, Typography } from '@mui/material';

export default function CardModal({ card, open, onClose }) {
  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>{card?.name}</DialogTitle>
      <DialogContent>
        <Typography variant="body1" gutterBottom>
          {card?.description}
        </Typography>
        <Typography variant="subtitle2" color="textSecondary">
          Ключевые слова: {card?.keywords}
        </Typography>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Закрыть</Button>
      </DialogActions>
    </Dialog>
  );
}