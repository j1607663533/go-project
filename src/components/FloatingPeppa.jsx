import React, { useState, useEffect, useRef } from 'react';
import './FloatingPeppa.css';

function FloatingPeppa() {
  const [position, setPosition] = useState({ x: 50, y: 50 });
  const [velocity, setVelocity] = useState({ vx: 2, vy: 2 });
  const [isDragging, setIsDragging] = useState(false);
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 });
  const containerRef = useRef(null);

  useEffect(() => {
    if (isDragging) return;

    const updatePosition = () => {
      if (containerRef.current) {
        const rect = containerRef.current.getBoundingClientRect();
        let { x, y } = position;
        let { vx, vy } = velocity;

        x += vx;
        y += vy;

        // Bounce off the right and left walls
        if (x + rect.width > window.innerWidth || x < 0) {
          vx = -vx;
        }

        // Bounce off the top and bottom walls
        if (y + rect.height > window.innerHeight || y < 0) {
          vy = -vy;
        }

        setPosition({ x, y });
        setVelocity({ vx, vy });
      }
    };

    const animationFrame = requestAnimationFrame(updatePosition);

    return () => cancelAnimationFrame(animationFrame);
  }, [position, velocity, isDragging]);

  const handleMouseDown = (e) => {
    setIsDragging(true);
    setDragStart({
      x: e.clientX - position.x,
      y: e.clientY - position.y,
    });
  };

  const handleMouseMove = (e) => {
    if (isDragging) {
      setPosition({
        x: e.clientX - dragStart.x,
        y: e.clientY - dragStart.y,
      });
    }
  };

  const handleMouseUp = () => {
    setIsDragging(false);
  };

  return (
    <div
      ref={containerRef}
      className="floating-peppa-container"
      style={{ top: `${position.y}px`, left: `${position.x}px` }}
      onMouseDown={handleMouseDown}
      onMouseMove={handleMouseMove}
      onMouseUp={handleMouseUp}
      onMouseLeave={handleMouseUp} // Also stop dragging if the mouse leaves the element
    >
      <div className="peppa-pig-circle">
        <span>Peppa Pig</span>
      </div>
    </div>
  );
}

export default FloatingPeppa;
