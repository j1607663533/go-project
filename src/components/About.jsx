import React from 'react';
import { Card, CardContent, Typography } from '@mui/material';

function About() {
  return (
    <Card>
      <CardContent>
        <Typography variant="h5" component="div">
          About
        </Typography>
        <Typography sx={{ mt: 2 }} color="text.secondary">
          This is the about page.
        </Typography>
        <Typography sx={{ mt: 1 }}>
          Here you can find more information about this application.
        </Typography>
      </CardContent>
    </Card>
  );
}

export default About;
