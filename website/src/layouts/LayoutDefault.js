import React from 'react';
import Footer from '../components/layout/Footer';
import Header from '../components/layout/Header';

const LayoutDefault = ({ children }) => (
  <>
    <Header />
    <main className="site-content">
      {children}
    </main>
    <Footer />
  </>
);

export default LayoutDefault; 