import { useEffect, useRef } from 'react';
import * as d3 from 'd3';

const TreeChart = ({ data }) => {
  const svgRef = useRef();

  useEffect(() => {
    if (!data) return;

    const margin = { top: 40, right: 120, bottom: 40, left: 120 };
    d3.select(svgRef.current).selectAll('*').remove();

    const svg = d3.select(svgRef.current).append('g');

    svg.append('defs')
      .append('marker')
      .attr('id', 'arrow')
      .attr('viewBox', '0 0 10 10')
      .attr('refX', 5)
      .attr('refY', 5)
      .attr('markerWidth', 6)
      .attr('markerHeight', 6)
      .attr('orient', 'auto-start-reverse')
      .append('path')
      .attr('d', 'M 0 0 L 10 5 L 0 10 z')
      .attr('fill', '#555');

    const root = d3.hierarchy(data, d => d.children);
    const treeLayout = d3.tree().nodeSize([120, 100]);
    treeLayout(root);

    const allX = root.descendants().map(d => d.x);
    const allY = root.descendants().map(d => d.y);
    const minX = Math.min(...allX), maxX = Math.max(...allX);
    const minY = Math.min(...allY), maxY = Math.max(...allY);
    const width = maxX - minX + margin.left + margin.right;
    const height = maxY - minY + margin.top + margin.bottom;

    d3.select(svgRef.current)
      .attr('width', width)
      .attr('height', height);

    svg.attr('transform', `translate(${-minX + margin.left}, ${-minY + margin.top})`);

    svg.selectAll('path.link')
      .data(root.links())
      .join('path')
      .attr('class', 'link')
      .attr('fill', 'none')
      .attr('stroke', '#555')
      .attr('stroke-width', 2)
      .attr('marker-start', 'url(#arrow)')
      .attr('d', d => `
        M${d.target.x},${d.target.y}
        V${(d.source.y + d.target.y) / 2}
        H${d.source.x}
        V${d.source.y}
      `);

    const node = svg.selectAll('g.node')
      .data(root.descendants())
      .join('g')
      .attr('class', 'node')
      .attr('transform', d => `translate(${d.x},${d.y})`);

    const boxWidth = 100;
    const boxHeight = 60;
    const iconSize = 24;

    node.append('rect')
      .attr('x', -boxWidth / 2)
      .attr('y', -boxHeight / 2)
      .attr('width', boxWidth)
      .attr('height', boxHeight)
      .attr('fill', '#f0f0f0')
      .attr('stroke', '#444')
      .attr('rx', 6)
      .attr('ry', 6);

    const baseUrl = 'http://localhost:8080/svgs/';
    const fallbackIcon = baseUrl + 'default.svg';

    // Fetch recipe.json
    fetch('http://localhost:8080/recipe.json')
      .then(res => res.json())
      .then(recipeData => {
        // Create a lookup map: { [element]: local_svg_path }
        const pathMap = {};
        recipeData.forEach(item => {
          if (item.element && item.local_svg_path) {
            pathMap[item.element] = item.local_svg_path;
          }
        });

        node.each(function (d) {
          const group = d3.select(this);
          const rawPath = pathMap[d.data.element];
          const imageUrl = rawPath ? baseUrl + rawPath.replace(/\\/g, '/') : fallbackIcon;

        group.append('image')
          .attr('x', -iconSize / 2)
          .attr('y', -boxHeight / 2 + 6)
          .attr('width', iconSize)
          .attr('height', iconSize)
          .attr('href', imageUrl)
          .on('error', () => {
            // Optional fallback if SVG fails to load
            d3.select(d3.event.target).attr('href', fallbackIcon);
          });
        });
      })
      .catch(err => {
        console.error('Failed to load recipe.json:', err);
      });

    node.append('text')
      .attr('text-anchor', 'middle')
      .attr('y', iconSize / 2 + 2)
      .attr('dy', '0.6em')
      .text(d => d.data.element)
      .style('font-size', '12px')
      .style('fill', '#111');

  }, [data]);


  return (
    <div style={{
      width: '100%',
      height: '100%',
      overflow: 'auto',
      border: '1px solid #ccc',
      borderRadius: '8px',
      background: 'white'
    }}>
      <svg ref={svgRef}></svg>
    </div>
  );
};

export default TreeChart;
