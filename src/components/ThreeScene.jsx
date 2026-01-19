import React, { useRef, useState } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { Box, useTexture, OrbitControls, DragControls } from '@react-three/drei';
import PeppaPig from './PeppaPig';

function TexturedCube() {
  const meshRef = useRef();

  // The order of textures corresponds to:
  // right, left, top, bottom, front, back
  const textures = useTexture([
    'https://placehold.co/128x128/ff0000/ff0000.png', // right
    'https://placehold.co/128x128/00ff00/00ff00.png', // left
    'https://placehold.co/128x128/0000ff/0000ff.png', // top
    'https://placehold.co/128x128/ffff00/ffff00.png', // bottom
    'https://placehold.co/128x128/800080/800080.png', // front
    'https://placehold.co/128x128/ffa500/ffa500.png', // back
  ]);

  useFrame(() => {
    if (meshRef.current) {
      meshRef.current.rotation.x += 0.005;
      meshRef.current.rotation.y += 0.005;
    }
  });

  return (
    <Box ref={meshRef} position={[-2, 0, 0]}>
      {textures.map((texture, index) => (
        <meshStandardMaterial key={index} attach={`material-${index}`} map={texture} />
      ))}
    </Box>
  );
}

function ThreeScene() {
    const [isDragging, setIsDragging] = useState(false);
    const orbitControlsRef = useRef();
  return (
    <div style={{ width: '100%', height: '500px' }}>
      <Canvas>
        <ambientLight intensity={0.5} />
        <pointLight position={[10, 10, 10]} />
        <DragControls onDragStart={() => setIsDragging(true)} onDragEnd={() => setIsDragging(false)}>
            <TexturedCube />
            <group position={[2, 0, 0]}>
                <PeppaPig />
            </group>
        </DragControls>
        <OrbitControls ref={orbitControlsRef} enabled={!isDragging} />
      </Canvas>
    </div>
  );
}

export default ThreeScene;
