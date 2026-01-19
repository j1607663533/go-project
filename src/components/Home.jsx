import React from 'react';
import { Card, CardContent, Typography } from '@mui/material';

function Home() {
  return (
    <Card>
      <CardContent>
        <Typography variant="h5" component="div">
          Home
        </Typography>
        <Typography sx={{ mt: 2 }} color="text.secondary">
          Welcome to the home page!
        </Typography>
        <Typography sx={{ mt: 1 }}>
          This is a simple React application with Vite, React Router, and Material-UI.
        </Typography>
      </CardContent>
    </Card>
  );
}

export default Home;
