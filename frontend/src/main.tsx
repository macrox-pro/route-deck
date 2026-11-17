import '@/styles/globals.css';

import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import { RouteProvider } from '@/providers/route-provider';
import { ThemeProvider } from '@/providers/theme-provider';

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <BrowserRouter>
            <RouteProvider>
                <ThemeProvider>
                    <h1>Hello World</h1>                    
                </ThemeProvider>
            </RouteProvider>
        </BrowserRouter>
    </React.StrictMode>
);