import React, { useRef } from 'react';
import { useFrame } from '@react-three/fiber';
import { Box } from '@react-three/drei';

function PeppaPig() {
  const groupRef = useRef();
  const pink = '#ffc0cb';
  const darkPink = '#ff69b4';

  useFrame(({ clock }) => {
    if (groupRef.current) {
      // Gentle bobbing motion
      groupRef.current.position.y = Math.sin(clock.getElapsedTime() * 2) * 0.1;
    }
  });

  return (
    <group ref={groupRef}>
      {/* Head */}
      <Box position={[0, 0.5, 0]}>
        <meshStandardMaterial color={pink} />
      </Box>
      {/* Snout */}
      <Box position={[0, 0.5, 0.5]} scale={[0.5, 0.3, 0.5]}>
        <meshStandardMaterial color={darkPink} />
      </Box>
       {/* Ears */}
       <Box position={[-0.3, 1, 0]} scale={[0.2, 0.4, 0.2]}>
        <meshStandardMaterial color={pink} />
      </Box>
      <Box position={[0.3, 1, 0]} scale={[0.2, 0.4, 0.2]}>
        <meshStandardMaterial color={pink} />
      </Box>
      {/* Body */}
      <Box position={[0, -0.75, 0]} scale={[1, 1.5, 1]}>
        <meshStandardMaterial color="red" />
      </Box>
      {/* Legs */}
      <Box position={[-0.2, -1.75, 0]} scale={[0.1, 0.5, 0.1]}>
        <meshStandardMaterial color="black" />
      </Box>
      <Box position={[0.2, -1.75, 0]} scale={[0.1, 0.5, 0.1]}>
        <meshStandardMaterial color="black" />
      </Box>
    </group>
  );
}

export default PeppaPig;
