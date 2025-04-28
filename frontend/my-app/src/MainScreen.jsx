import { Button, Grid, Typography } from '@mui/material';
import { useNavigate } from 'react-router-dom';

export default function MainScreen() {
  const navigate = useNavigate();

  return (
    <div style={{ padding: 20 }}>
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <Button 
            variant="contained" 
            fullWidth
            onClick={() => navigate('/daily')}
          >
            Ежедневное гадание
          </Button>
        </Grid>
        
        <Grid item xs={12}>
          <Button 
            variant="outlined" 
            fullWidth
            onClick={() => navigate('/ask-question')}
          >
            Задать вопрос
          </Button>
        </Grid>
      </Grid>
    </div>
  );
}