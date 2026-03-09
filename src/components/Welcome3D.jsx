import React, { useRef, useMemo } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { OrbitControls, Float, Center, Sparkles, ContactShadows, Environment } from '@react-three/drei';
import * as THREE from 'three';

// 100% 稳定的方案：将文字绘制到 Canvas 贴图上，不依赖外部字体文件
const Text3DStable = ({ text }) => {
  const meshRef = useRef();
  const layers = 10; // 更多层数产生更强的立体感
  const offset = 0.015;

  // 创建一个包含文字的贴图
  const texture = useMemo(() => {
    const canvas = document.createElement('canvas');
    canvas.width = 1024;
    canvas.height = 256;
    const ctx = canvas.getContext('2d');
    
    // 背景透明
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    
    // 设置文字样式
    ctx.font = 'bold 160px "Microsoft YaHei", "PingFang SC", sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    
    // 文字渐变色
    const gradient = ctx.createLinearGradient(0, 0, canvas.width, 0);
    gradient.addColorStop(0, '#667eea');
    gradient.addColorStop(1, '#764ba2');
    
    // 绘制文字
    ctx.fillStyle = gradient;
    ctx.fillText(text, canvas.width / 2, canvas.height / 2);
    
    const tex = new THREE.CanvasTexture(canvas);
    tex.needsUpdate = true;
    return tex;
  }, [text]);

  useFrame((state) => {
    if (meshRef.current) {
      const t = state.clock.elapsedTime;
      // 稍微晃晃，增加仪式感
      meshRef.current.rotation.y = Math.sin(t * 0.4) * 0.15;
    }
  });

  return (
    <Center>
      <Float speed={2} rotationIntensity={0.5} floatIntensity={0.5}>
        <group ref={meshRef}>
          {[...Array(layers)].map((_, i) => (
            <mesh key={i} position={[0, 0, -i * offset]}>
              <planeGeometry args={[5, 1.25]} />
              <meshStandardMaterial 
                map={texture} 
                transparent={true} 
                opacity={i === 0 ? 1 : 0.8} // 后面几层降低点透明度产生深度感
                side={THREE.DoubleSide}
                metalness={0.6}
                roughness={0.4}
                depthWrite={i === 0} // 只有第一层写深度，防止透明排序异常
              />
            </mesh>
          ))}
        </group>
      </Float>
    </Center>
  );
};

const Welcome3D = () => {
  return (
    <div style={{ 
      width: '100%', 
      height: '180px', 
      cursor: 'grab',
      margin: '0 auto',
      position: 'relative',
      overflow: 'hidden',
    }}>
      <Canvas 
        shadows 
        camera={{ position: [0, 0, 4], fov: 40 }}
        gl={{ antialias: true, alpha: true }}
      >
        <ambientLight intensity={1} />
        <pointLight position={[10, 10, 10]} intensity={1.5} />
        
        {/* 不再使用 Suspense 包裹 Text，因为我们不再加载外部资源 */}
        <Environment preset="city" />
        <Text3DStable text="欢迎登录" />
        <Sparkles count={40} scale={6} size={1.2} speed={0.5} color="#667eea" />
        <ContactShadows position={[0, -0.8, 0]} opacity={0.4} scale={10} blur={2.5} />

        <OrbitControls 
          enableZoom={false} 
          enablePan={false}
          minPolarAngle={Math.PI / 3}
          maxPolarAngle={Math.PI / 1.5}
        />
      </Canvas>
    </div>
  );
};

export default Welcome3D;
