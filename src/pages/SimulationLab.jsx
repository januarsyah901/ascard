import React, { useState, useCallback, useEffect } from 'react';
import { Link } from 'react-router-dom';

const GRID_IMAGES = [
  "https://lh3.googleusercontent.com/aida-public/AB6AXuCK1oa99_FjjSkbsHXgmNz3FZSpAXn_gkQdDyJlfQn6urldOQMga25FZQYotbDJdDZVIaLlTG4mjvZ4w5OPGkwNmFzO6m2Ct11Sb-2UBZhLNvv2GthNisvwXB_oaUTqaDbNyxMu6tr8loe0E2J9T1E6y-GfvCfnr9Fn8tHe6iS7scpB4cRupX8nxSxr72rGfh8lAOS1AWzp-bx7dKAn16_i1zmX0mwuB0t4t08kYyt0JdOAi8lhfz_TpyjsV7CFcQ0OmQGNFFc5Ukh6",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuAizptxBhzGJOGfUgJEo71V3s-xlskuR8XNBNW_cfFMEkVk-Ke-0q9lpHVPMb2KmoBX0ZhpAPEh09xRELnG0_L9S-zju5rirQrjv-qq1Bzx7gXPHE49tPf9h2DJ5ExEt_LFxzLQGmjDc7MXfD7FghdaKYIKwdmf83bgfLBuUGi8vhE9dNqx5Zn18JQzsp_Fiqfgvig9I0JDTdgchANScvXb4j8NWHAgffOcdClVdU39NFJCHOkqv1jz7iCk6g5j6mTijtpFO3eQ-9Ze",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuDu1ottpgEoXzBjHYpscgM9uDHYreDhCbCOZYqaU3QpDiQudoYFKX-FshHEIVtUsROvntxNUO1Al7E-kwgVjCvkrjgUNLG_GJwZ310PlRcr2SluR2gBxrIAI0_MgqJ0EOxil2oatNZgTILkBoRLgaiXlQSq-2uhkBA9GCaHFGF9uwMKypt2ExybU899T1yt_grmCAS7tmQc00pjNJHEpbvf85OHixEqJxJfIsvymhqnkWuS-xiRJ0Rx_vBJgE_9doagnm-a5E7HJr5E",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuB0RbXhInkkoYRHViCWPStGlZBQVHSkJquioeSmzcvPyanjemCQYFrApC8OtANA5yAx2F4eUxRu-lx5j30EfMf_8tGvKOFvPMJv27dVGZ7O76Rb0ndouKnf7_rTNPcpP7iFxnfdt-WjEd8OiKUcGKr69kLKnjG484SNXvZzwvpt6a8AOAXnobZAInif7aV3yQh80Je7l-Vu0J-6WduRAZlW389YseRWmeNgh19dFAmsrHMqAm6r5D-8ckwZ-pfmHfWJrYeRiIMbXMKN",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuCESJCAMXjVzVu2TDLQyHUO3yu91-glEciCBOjaI9c1Sc6qkC9lSUgAH_-mloQmFYQUuH298TsBjWHYWDgjZrgmrQDEwKNL4wtNP3hJFrNMOvwuALVC-leBgYoUf8qRWpqVxI5U_suo0NVtcNPdfD5VE_7v8xZVpuuII__yTmlu7btLDQJGmL4stsDTD44mJOeGyn56zBVLyVAZrKgeAAcXV15yzLp3dVkm1c7IY4cn5ZfTGrAhvkFxjrTj8_MasK6ATe-ZOHUpIJ6M",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuC-F_2oZmz6Cl3rf1oxXj17dBPe5STdWFJPCTB64j7oDYZ3KBgVh8vUJVR5xcIFf_GUbeF9qysvFPvtH1AsqgxpskpxbPrzPQDIHRsvuEcetfTTRk7IuDPKlZoQXtppeD4yfQTyt7at9SQMc9FppDT4v9GG7BMQXaNQfnWXbFqCUZG-HVRS4ZSS09qnkpmY8uu20KRrL9MZPZOB9cn6PKeVpBqTRzcVd9mWmEIwtre9E7F1SSOHW4VzhVzgc_QO_Dc1Ntt0hrJI-5ua",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuCafkEM33Q6ijEqayDxBgEkh1douBAHjBY5dXNH4juDn7N35IIQAx2vRmlt6W77dmp9pI2ypLUCq9PRz8p09yJn1FKu7c3IOiZB38Zv_vB_yx8bH8VzX2zBfGM_tj6xbYlK-V6MdHNCWyQjaJwOrdxrIySWgycluWDhhcOSBb2DPPoojRG_ImeHaRTQvs2ORXqred6WUSPcI-ib15LrIiYg899D1JVcCE8YCaFjgh3Mp_40D9UMaZKKhojKm7Fs3riFdggdihWRzTvY",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuCIlNhnmr7EN0m6aESwVNfuIEpa0yGyAahLQyNwHV5k9xo1yI3kanR5hIWkSn_hzUvJ3j6MIVwGcJxk7GeLJPRq8tKgZ-9N_s3NOfwOqdJLHY1E0GTENjCVpALBzheUNItf9CeUGHLA--B3Ip9Te5K84KYXv76nqjwFftwW_XnOiju2c0eTeqOTMIbfZ85xGd2dKKSnYXx_Sa-i-w5Q5iYb2vskPhZY1O4t04I5L-2GMB1y4Q8IkR6GuNbDRsd33XiH7n8LVuHWtMiW",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuDtkyqyR5ukjv6NZu70TaMCODTD4ooJDoaJ_9bkUOAxSMsGSnf7WAu68PO9y3S-31M5_RdVgbgtEp2w_EOST_ATskINq_elAoF6ds_Sd71RO4MUTEjfGiEtPlMASPdTJWqdsw686W70xEcceWZ0kg2cEdukKw2tozaDQ7B0fqNs0asIK6SaU-wTcd-RZeLj_3mF6N0nXaK2qcziR_pv2mdGcfsJwLW-Di-ymieop4ngaX3LBL-RqF9H4gvI5LQdxtqikCu3Md09Khrf",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuAHlQTWNJ6npnCpAMeSWpx9uN-bIhdqjr6VdQGcQMhlzfV6NKQy53qeZ-Uynvo-OoX7dwKV3v7wL8WI0E9S68CLf9UeP7mLpLjZHFgV7tu7In6nQUdgY195uuuxnDJ6K5RkisjvVMor7mb72rUAS5DhwEtyNBRmBVgM6K6c-zb8gLvUVmeq1wibWPArOhU67A0drmn_bjjdTm1IT4nQDnhLvd8Tau8ErANketCNzLIDjUqc3FOmVB6reyolubzxjgQMvA-Eydl02fOx",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuAR7plXW3UnLjgnhGr3TR-HTLExtsY3bq_76hjrSSDJKj0d6Ge0O-d1wwOEAtAJE_39kfkVlOBcwlkFxzzSYgcWoujLt1S6a4k2IHA5T6p0jw6jFKv_9EtgBXK820MCbUOU1terXxaJis_S89bA9D_LwoZiVYKfScPkiko7zWme1kX9m17NDZii-D-wT_23i8rE5vaVNF3vYWikBZGPOdTUYNg1sYWkIYyNmx3TGbC3nV_hZSIJ4a3E4verhF1dAL1UDKhXTegAt7gQ",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuC5VO5u1VAqShYyzWaTg1lroBFSQKlH2AqJyphDy3ouX9NqvU-V_6Lv8gdksKAchEV3NWMs4NJnZ1hnOg7nWBY35rEQcAItjEElzkyObdxLr75ozCd3-s8dblCAb38i6S5o9SRNBoRtAWh3MF8N9dXQpINBQV3fVB_KwBTKAlTVcbauQGqfmMEkBqiDblzFtJAZn5Wq4r98nBpnuDNaP-Wbxbm0zg9vNmUMwMVZARIUUa4aZ8AXGJqt5cpbntU6HrViCn6u7eZlpbNI",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuB1K5CkrqZI621F0OxnbzvaM0tndHsjifZ9ilQ05odFcI_-o1AK_mFtXn9mWJqF9D8fd_RrVbhnvXhJ1idQW-D1xiyZ-awou-9HvDXuXwBKztmmib8_PcbkEJOhXept-5v4JKDml637RwRbG1UiXEmVESxf9OYLEW-Qtmln7jRwphDVDypXjgRtma_grEaUkzPkw2WNAPZxHYv_i0xeBAj2KN5o1I3Pd4JkxeHXTjtCZTp6G4i6JEYv7elx7JUPfOvSJn8EDUaTYBgN",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuCL-aL-FScpmu_E4PG0ep59T9CxDnpqlV_IAOrsCyF-cvVSZLyV5OMXNaHbkXQN59LOrVD-TDDcCxlh0nWZ1k4Qxva_yCoaoceshEMBYdrO1o0wnoebb6jl4fVpGITwTKeRgZ7OJ4onFLH4z8BFNE9uoL9deuGrIR9DVnSEkcWySr7WVg8Botmk3r1k5wlYDjeg0-C4pF6fwVjpTyS5UmsaOfeXLCRMdClbhy8mIjq8JffG2bIkOB3x-edfZozyJAtddh4KtYs7miWD",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuCt6yt0VeVBU_yF4xT4Hg3hk0Nv1DcVSVqlV57B7_02tCK7cqqni0BzlrUj7hAQSPA78osqCbw4Awn1gnj4_UHipp69i7qT7FigCRMnZuXW0AXk_MHt19u6AcdRaHY1r1m8gkZwavK_UtJJTBUaXvS1BJPLGDD8wvZm_pKafPUhqRJ-UzCbKsE7lpw_RcY2ZGpy09y9w45DtW7juWwS9_zoBn9cEDY6YgTJEiAzHM0i7Ftz_wwMfg9oC80bPmjwgRcqTJ9K3bknT0Ja",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuCAaSu5gWGruz0S7JTuZ8hJPj8HnV2zfQI6tXcPPezb2lBDOJc7keRfQtbjMkEDj9EWZLDMcwXasYzt6eOPbjlWHasSlfMnTYhT9BRaJ3FkaMy4wKDCy3K3c48bNpY_xJPLDiwLMQFO0SWji7PrxActDjj8XOecivRisKQ_SzWB4uoj3GZrGgJVfjex04bEer7YuuJFeHOxUHK_Nqg9alr5ecP-I92XgL6Y3XABzb7e5CKA91AXUat0tZxT6T3DmNhR-pjbDS7VFpg2",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuAL0qefcpE-rmgxZ0a-7P0GvdKNP0SzG3GGAztCVyUZxhfgJrJxInSCQ04uUMZcU7n_WIbhVqPo-beIJS15drjOGrf3sqgzPTY0Hm5L6fCiB3wzIt7BXKiEfaXd5VSel8JLyqlM0uhKSr02JdzorP2nbyYVw5cBNJec69Bh7d0Bn88YZv-5r5mjCTUtinGvqob7IMwDfLdHVl0tywqZcfqFKP_V-YPKfQ5w24UZljfn72fbz6oyI24tnrrSzIK_DBkLSoKRZQvM5YFv",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuC6z4z2tVB1Zxf9g3HPYiTsxhIRe7N8wTLsTv92elyAttbwzu01rPgBqUM561d6h6DR1EBKdzBNA9hwzxV4Zqw5a5mEXxxUAoyxnHOeHu905dqk7TjbYhUInAxsnalr0t9fVUArCBV05ORrP2PUG_kaqxMtv1kl_2A4YecAysARvsfdvZ4uAlckT7x7K3whmcllFTFS9RCB34rCEIp9Df4GW_m4rhb5uZNV5nwlCp08ygAwV6hmCIxLODSgFFZeYaGvxg-tu0EzWMGK",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuAghkP4QB3Xw07wj8BRVar4znczIRwvZlwXIDCzf6JBYgXs-xsibxSDt6DrPdJfdnMpoFAYoY0TdP8_TXPLQnh6DWADRe9TzvN565j402KE_nsQuKCqJBiIMAKYJeAPzaTD3HEi-Us4Je15oExCSQZDQACz01o61ZpKoLV9ApwNox1macZhw1WXe3uLbKv4kx1ysXE0TygxxB4tIrUx3cjGVsr6PQuYUXnm0abzmYAqZqcC8KIY05AOSBJCe9zr6UbABLCJcGMRewhK",
  "https://lh3.googleusercontent.com/aida-public/AB6AXuDg2tcbC_FvRBa_PVyN1OhDO3fIt3B2aYXinRRHm4_-FCszHIyE9HXCHyYFiyB7i-o1BFEw9LJ_ZE47_tygABeNOS9rs_fmtXJ4F5xoed40g6SFhYjS0ijlyuuYseC_97WepmZwxvY7N60zsXLjiJhnTGH2QKRu9y6OVTg_PiGkL4oWp0Y1N3_3Q5m3RdwnXTBtx75gxOKaD8M_Vl4PcavONnT39qeIvsIAdzT5hiM8GTDwbR4eF59n_OY8V7bOzcpKTeUDzWTpeHGp"
];

import { useLanguage } from '../context/LanguageContext';
import { useCurrency } from '../context/CurrencyContext';

const IDR_RATE = 15000;
const INITIAL_BALANCE_USD = 5000000 / IDR_RATE;
const DEFAULT_BID_USD = 10000 / IDR_RATE;

export default function SimulationLab() {
  const { t } = useLanguage();
  const { formatValue, currency } = useCurrency();
  const [balance, setBalance] = useState(INITIAL_BALANCE_USD);
  const [totalLoss, setTotalLoss] = useState(0);
  const [spinId, setSpinId] = useState(1);
  const [currentBid, setCurrentBid] = useState(DEFAULT_BID_USD);
  const [isCustomBid, setIsCustomBid] = useState(false);
  const [customBidInput, setCustomBidInput] = useState('');

  const [logs, setLogs] = useState([]);
  
  const [rngData, setRngData] = useState([
    `> [${Math.random().toFixed(6)}, ${Math.random().toFixed(6)}, ${Math.random().toFixed(6)}]`,
    `> HASH: 8f2a...f92c (READY)`,
  ]);
  
  const [isSpinning, setIsSpinning] = useState(false);
  const [consecutiveLosses, setConsecutiveLosses] = useState(0);
  
  // Scramble grid items for visual effect
  const [gridImages, setGridImages] = useState(GRID_IMAGES);
  
  // New States for Suggestions 1 & 5
  const [history, setHistory] = useState([INITIAL_BALANCE_USD, INITIAL_BALANCE_USD, INITIAL_BALANCE_USD, INITIAL_BALANCE_USD]); // History of balance for the chart
  const [sessionSpins, setSessionSpins] = useState(0);
  const [showRealityCheck, setShowRealityCheck] = useState(false);

  // Suggestion 1: Simple SVG Line Chart Component
  const SimpleLineChart = ({ data }) => {
    if (!data || data.length < 2) return null;
    const min = Math.min(...data);
    const max = Math.max(...data);
    const range = max - min || 1;
    const width = 300;
    const height = 100;
    const points = data.map((val, i) => {
      const x = (i / (data.length - 1)) * width;
      const y = height - ((val - min) / range) * height;
      return `${x},${y}`;
    }).join(' ');

    return (
      <svg viewBox={`0 0 ${width} ${height}`} className="w-full h-full overflow-visible">
        <polyline
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          points={points}
          className="text-primary transition-all duration-500"
        />
        {/* Shadow area */}
        <polyline
          fill="currentColor"
          fillOpacity="0.1"
          stroke="none"
          points={`0,${height} ${points} ${width},${height}`}
          className="text-primary transition-all duration-500"
        />
      </svg>
    );
  };

  const handleSpin = useCallback(() => {
    if (isSpinning || balance < currentBid) return;
    
    setIsSpinning(true);
    setBalance(prev => prev - currentBid);
    
    // Rapidly scramble grid and RNG to simulate "spinning" effect
    const spinInterval = setInterval(() => {
      setGridImages(prev => [...prev].sort(() => 0.5 - Math.random()));
      setRngData(prev => [
        `> [${Math.random().toFixed(6)}, ${Math.random().toFixed(6)}, ${Math.random().toFixed(6)}]`,
        ...prev
      ].slice(0, 6));
    }, 80);
    
    // Simulate spinning delay
    setTimeout(() => {
      clearInterval(spinInterval);
      
      const isWin = Math.random() < 0.04; // 4% hit rate to emphasize loss
      let winAmount = 0;
      
      if (isWin) {
        // Most "wins" are actually losses (e.g. bet $2, win $0.40)
        winAmount = Math.random() < 0.8 ? (Math.random() * 1.5) : (Math.random() * 10);
      }
      
      const netChange = isWin ? winAmount : -currentBid;
      
      if (isWin && winAmount >= currentBid) {
        setConsecutiveLosses(0);
      } else {
        setConsecutiveLosses(prev => prev + 1);
      }
      
      setBalance(prev => {
        const nextBalance = prev + winAmount;
        setHistory(h => [...h, nextBalance].slice(-20)); // Keep last 20 points
        return nextBalance;
      });
      setTotalLoss(prev => prev + netChange);
      
      setSessionSpins(s => {
        const next = s + 1;
        if (next % 15 === 0) {
          setShowRealityCheck(true);
        }
        return next;
      });

      setSpinId(prevId => {
        const nextId = prevId + 1;
        setLogs(prevLogs => [
          {
            id: nextId,
            type: isWin ? 'win' : 'loss',
            amount: isWin ? winAmount : netChange
          },
          ...prevLogs
        ].slice(0, 3));
        return nextId;
      });
      
      // Update RNG Data stream final state
      setRngData([
        `> [${Math.random().toFixed(6)}, ${Math.random().toFixed(6)}, ${Math.random().toFixed(6)}]`,
        `> [${Math.random().toFixed(6)}, ${Math.random().toFixed(6)}, ${Math.random().toFixed(6)}]`,
        `> HASH: ${Math.random().toString(16).substring(2, 6)}...${Math.random().toString(16).substring(2, 6)} (UNVERIFIED)`,
        `> [${Math.random().toFixed(6)}, ${Math.random().toFixed(6)}, ${Math.random().toFixed(6)}]`,
        `> [${Math.random().toFixed(6)}, ${Math.random().toFixed(6)}, ${Math.random().toFixed(6)}]`,
        `> [${Math.random().toFixed(6)}, ${Math.random().toFixed(6)}, ${Math.random().toFixed(6)}]`
      ]);
      
      // Shuffle grid images randomly one last time
      setGridImages([...GRID_IMAGES].sort(() => 0.5 - Math.random()));
      
      setIsSpinning(false);
    }, 800); // 800ms fake delay
    
  }, [balance, isSpinning]);

  const handleTopup = (amount) => {
    setBalance(prev => prev + amount);
  };

  const handleReset = () => {
    setBalance(INITIAL_BALANCE_USD);
    setTotalLoss(0);
    setSpinId(1);
    setConsecutiveLosses(0);
    setLogs([]);
    setCurrentBid(DEFAULT_BID_USD);
    setIsCustomBid(false);
  };

  let tiltLevel = 'Low';
  let tiltColor = 'bg-black w-[25%]';
  let analystMessage = "Observation: Your clicking cadence is steady. The RNG continues to generate true random outcomes with a fixed 96.4% RTP.";
  
  if (consecutiveLosses > 15) {
    tiltLevel = 'Critical';
    tiltColor = 'bg-error w-[95%]';
    analystMessage = "Warning: You've experienced 15+ consecutive sub-bet returns. The 'Gambler's Fallacy' might convince you a win is due. It is not. The system has no memory.";
  } else if (consecutiveLosses > 8) {
    tiltLevel = 'High';
    tiltColor = 'bg-black w-[75%]';
    analystMessage = "Observation: Your click frequency is likely increasing after this streak of losses. This is a common cognitive bias called 'Chasing.' Mathematically, each spin is a fresh 0.05% probability.";
  } else if (consecutiveLosses > 4) {
    tiltLevel = 'Medium';
    tiltColor = 'bg-black w-[50%]';
    analystMessage = "Note: Experiencing sequential losses is mathematically expected in a high-volatility environment. Do not interpret this as a 'Cold' streak that will soon become 'Hot'.";
  }

  return (
    <div className="flex bg-surface" style={{ minHeight: '100vh' }}>
      {/* SideNavBar (Left Fixed) */}
      <aside className="hidden lg:flex flex-col w-64 bg-neutral-50 border-r border-outline-variant/30 overflow-y-auto shrink-0">
        <div className="p-6">
          <Link to="/" className="font-bold text-black text-xs tracking-tighter block hover:underline">BONGKAR ILUSI</Link>
          <span className="text-[9px] text-neutral-400 font-medium uppercase tracking-widest block mt-1">DATA LAB V1.0.4</span>
        </div>
        <nav className="flex-grow flex flex-col mt-4">
          <Link to="/simulation" className="bg-neutral-200 text-black font-bold px-6 py-4 transition-colors duration-150 flex items-center gap-3 text-[11px] uppercase tracking-tight">
            <span className="material-symbols-outlined text-sm">analytics</span> {t.nav.simulator}
          </Link>
          <Link to="/intervention" className="text-neutral-500 px-6 py-4 hover:bg-neutral-100 transition-colors duration-150 flex items-center gap-3 text-[11px] uppercase tracking-tight">
            <span className="material-symbols-outlined text-sm">science</span> {t.nav.logic}
          </Link>
          <Link to="/knowledge-base" className="text-neutral-500 px-6 py-4 hover:bg-neutral-100 transition-colors duration-150 flex items-center gap-3 text-[11px] uppercase tracking-tight">
            <span className="material-symbols-outlined text-sm">visibility</span> {t.nav.library}
          </Link>
        </nav>
        <div className="p-6 space-y-6">
          <button onClick={handleReset} className="w-full bg-black hover:opacity-80 text-white text-[10px] font-bold tracking-widest py-4 uppercase transition-colors">{t.simulation.reset}</button>
          <div className="space-y-4">
            <div className="flex items-center gap-3 text-[11px] uppercase tracking-tight text-neutral-500 cursor-pointer hover:text-black">
              <span className="material-symbols-outlined text-sm">terminal</span> System Logs
            </div>
          </div>
        </div>
      </aside>
      
      {/* Main Workspace (Simulation Area) */}
      <main className="flex-grow bg-surface overflow-y-auto p-4 md:p-8 flex flex-col xl:flex-row gap-8 min-h-screen">
        {/* Area Simulasi */}
        <section className="flex-grow space-y-8 xl:max-w-4xl">
          {/* Header Mobile / Title */}
          <div className="flex justify-between items-end flex-wrap gap-4">
            <div>
              <div className="lg:hidden mb-4">
                <Link to="/" className="font-bold text-black text-xs tracking-tighter uppercase underline">Bongkar Ilusi</Link>
              </div>
              <h2 className="text-3xl sm:text-4xl font-black tracking-tighter uppercase Inter leading-none">{t.simulation.title}</h2>
              <p className="text-xs uppercase tracking-widest text-neutral-500 font-bold mt-2">{t.simulation.activeSession} #{spinId}</p>
            </div>
            <div className="text-right">
              <p className="text-[10px] uppercase font-bold tracking-widest text-neutral-400">{t.simulation.totalLoss}</p>
              <p className="text-2xl font-bold tracking-tighter text-error">{formatValue(totalLoss)}</p>
            </div>
          </div>
          
          {/* Simulation Grid (Mahjong Style - Desaturated) */}
          <div className={`grid grid-cols-5 gap-2 aspect-[5/4] bg-surface-container-low p-4 shadow-inner transition-opacity duration-150 ${isSpinning ? 'opacity-50' : 'opacity-100'}`}>
            {gridImages.map((src, idx) => (
              <div key={idx} className={`bg-white flex items-center justify-center p-4 ${idx % 3 === 0 ? 'border-2 border-black/10' : ''}`}>
                <img alt={`Grid item ${idx}`} className={`grayscale ${idx % 2 === 0 ? 'opacity-40' : 'opacity-70'}`} src={src} />
              </div>
            ))}
          </div>
          
          {/* Bid Selection */}
          <div className="bg-surface-container-low p-6 border-b border-black/5">
            <p className="text-[10px] font-black tracking-widest text-neutral-400 uppercase mb-4">Select Simulation Bid</p>
            <div className="flex flex-wrap gap-3">
              {[10000, 20000, 50000, 100000].map(bidIdr => {
                const bidUsd = bidIdr / IDR_RATE;
                const isActive = !isCustomBid && Math.abs(currentBid - bidUsd) < 0.01;
                return (
                  <button 
                    key={bidIdr}
                    onClick={() => { setCurrentBid(bidUsd); setIsCustomBid(false); }}
                    className={`px-4 py-2 text-[11px] font-black tracking-tight border-2 transition-all ${isActive ? 'bg-black text-white border-black' : 'border-black/10 hover:border-black text-black'}`}
                  >
                    {formatValue(bidUsd)}
                  </button>
                );
              })}
              <button 
                onClick={() => setIsCustomBid(true)}
                className={`px-4 py-2 text-[11px] font-black tracking-tight border-2 transition-all ${isCustomBid ? 'bg-black text-white border-black' : 'border-black/10 hover:border-black text-black'}`}
              >
                CUSTOM
              </button>
            </div>
            {isCustomBid && (
              <div className="mt-4 flex items-center gap-4 animate-in fade-in slide-in-from-top-2 duration-300">
                <div className="relative">
                  <span className="absolute left-0 top-1/2 -translate-y-1/2 text-xs font-bold text-neutral-400">{currency === 'IDR' ? 'Rp' : '$'}</span>
                  <input 
                    type="number"
                    placeholder={`Amount in ${currency}`}
                    value={customBidInput}
                    onChange={(e) => setCustomBidInput(e.target.value)}
                    className="bg-transparent border-b-2 border-black pl-6 pr-2 py-2 text-sm font-black focus:outline-none w-40"
                  />
                </div>
                <button 
                  onClick={() => {
                    const val = parseFloat(customBidInput);
                    if (val > 0) {
                      const usdVal = currency === 'IDR' ? val / IDR_RATE : val;
                      setCurrentBid(usdVal);
                    }
                  }}
                  className="bg-black text-white text-[10px] px-6 py-2 font-black uppercase tracking-widest hover:opacity-80 transition-opacity"
                >
                  Apply Bid
                </button>
              </div>
            )}
          </div>

          {/* Simulation Controls */}
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-6 h-auto sm:h-[132px]">
            <div className="bg-surface-container-low p-6 flex flex-col justify-center h-full">
              <p className="text-[10px] font-bold tracking-widest text-neutral-400 uppercase mb-1">{t.simulation.virtualBalance}</p>
              <p className="text-3xl font-black tracking-tighter">{formatValue(balance)}</p>
              <div className="mt-4 flex gap-2">
                <button onClick={() => handleTopup(1000000 / IDR_RATE)} className="bg-neutral-200 hover:bg-neutral-300 transition-colors text-black text-[10px] font-bold px-3 py-1 uppercase">{t.simulation.topup} {formatValue(1000000 / IDR_RATE)}</button>
              </div>
            </div>
            <div className="bg-primary p-6 flex flex-col justify-center items-center h-full">
              <button 
                onClick={handleSpin}
                disabled={isSpinning || balance < currentBid}
                className="w-full text-white font-black text-2xl tracking-tighter uppercase mb-1 flex items-center justify-center gap-2 hover:opacity-80 transition-opacity disabled:opacity-50"
              >
                <span className={`material-symbols-outlined ${isSpinning ? 'animate-spin' : ''}`}>refresh</span> 
                {isSpinning ? t.simulation.calculating : t.simulation.spin}
              </button>
              <p className="text-[10px] text-on-primary-fixed-variant tracking-widest uppercase mt-2">{t.simulation.cost.replace('{val}', formatValue(currentBid))}</p>
            </div>
            <div className="bg-surface-container-low p-6 overflow-hidden h-full flex flex-col">
              <p className="text-[10px] font-bold tracking-widest text-neutral-400 uppercase mb-3">{t.simulation.recentLogs}</p>
              <div className="space-y-2 text-[10px] font-medium font-mono text-neutral-600 flex-grow overflow-y-auto no-scrollbar">
                {logs.map((log) => (
                  <div key={log.id} className="flex justify-between border-b border-black/5 pb-1">
                    <span>#{log.id}</span>
                    <span className={log.type === 'loss' ? 'text-error' : 'text-black font-bold'}>
                      {log.type.toUpperCase()} ({log.amount > 0 ? '+' : ''}{formatValue(log.amount)})
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </section>
        
        {/* Right Sidebar: Transparency & AI Educator */}
        <aside className="w-full xl:w-80 2xl:w-96 flex flex-col gap-6 shrink-0">
          {/* Suggestion 1: Balance Analytics Chart */}
          <div className="bg-white p-6 border-b-2 border-black">
            <div className="flex justify-between items-center mb-4">
              <p className="text-[10px] font-bold tracking-widest text-neutral-400 uppercase">Balance Analytics (20 Iterations)</p>
              <span className="material-symbols-outlined text-sm text-primary">show_chart</span>
            </div>
            <div className="h-24 w-full">
              <SimpleLineChart data={history} />
            </div>
            <p className="text-[9px] text-neutral-400 uppercase mt-2 text-right tracking-widest">Volatility Model: High</p>
          </div>

          {/* Raw RNG Data (Terminal Style) */}
          <div className="bg-black p-6 flex flex-col h-48">
            <div className="flex justify-between items-center mb-4">
              <span className="text-[9px] text-white font-bold tracking-widest uppercase flex items-center gap-2">
                <span className="w-1.5 h-1.5 bg-white animate-pulse"></span> {t.simulation.rngStream}
              </span>
              <span className="text-[9px] text-neutral-500 uppercase">ms: {isSpinning ? '...' : '0.002'}</span>
            </div>
            <div className="flex-grow overflow-hidden font-mono text-[11px] text-neutral-400 space-y-1 no-scrollbar">
              {rngData.map((line, i) => (
                <p key={i} className={line.includes('HASH') ? 'text-white' : ''}>{line}</p>
              ))}
            </div>
          </div>
          
          {/* Real-time Odds Blocks */}
          <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-1 gap-4">
            <div className="bg-white p-6 border-b-2 border-black h-full">
              <p className="text-[10px] font-bold tracking-widest text-neutral-400 uppercase mb-2">{t.simulation.odds}</p>
              <div className="flex items-baseline gap-2">
                <span className="text-4xl font-black tracking-tighter">0.05%</span>
                <span className="text-[10px] text-neutral-400 uppercase font-bold">(1 IN 2000)</span>
              </div>
            </div>
            <div className="bg-white p-6 border-b-2 border-black h-full">
              <p className="text-[10px] font-bold tracking-widest text-neutral-400 uppercase mb-2">{t.simulation.rtp}</p>
              <div className="flex items-baseline gap-2">
                <span className="text-4xl font-black tracking-tighter">96.4%</span>
                <span className="text-[10px] text-neutral-400 uppercase font-bold">(CALCULATED)</span>
              </div>
            </div>
          </div>
          
          {/* AI Educator / Chat Bubble */}
          <div className="bg-surface-container-highest p-6 flex-grow flex flex-col">
            <div className="flex items-center gap-3 mb-6">
              <div className="w-10 h-10 bg-black flex items-center justify-center">
                <span className="material-symbols-outlined text-white text-xl">psychology</span>
              </div>
              <div>
                <p className="text-xs font-black tracking-tighter uppercase">{t.simulation.aiAnalyst}</p>
                <p className="text-[9px] text-neutral-500 uppercase tracking-widest">{t.simulation.activeMonitoring}</p>
              </div>
            </div>
            <div className="bg-white p-4 relative mb-6">
              <p className="text-xs leading-relaxed text-on-surface font-medium italic">
                "{analystMessage}"
              </p>
              <div className="absolute -bottom-2 left-6 w-4 h-4 bg-white rotate-45"></div>
            </div>
            <div className="mt-auto space-y-4">
              <div className="flex items-center gap-2">
                <div className="h-1 flex-grow bg-neutral-300">
                  <div className={`h-full ${tiltColor} transition-all duration-500`}></div>
                </div>
                <span className="text-[10px] font-bold uppercase">{t.simulation.tiltLevel}: {tiltLevel}</span>
              </div>
              <Link to="/intervention" className="block w-full border-2 border-black py-3 text-[10px] font-black text-center uppercase tracking-widest hover:bg-black hover:text-white transition-all">{t.simulation.break}</Link>
            </div>
          </div>
        </aside>
      </main>

      {/* Suggestion 5: Reality Check Modal */}
      {showRealityCheck && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/90 backdrop-blur-sm p-6">
          <div className="bg-white border-4 border-primary p-10 max-w-lg w-full shadow-[20px_20px_0px_0px_rgba(0,0,0,0.3)]">
            <div className="flex items-center gap-4 mb-8">
              <div className="w-12 h-12 bg-primary flex items-center justify-center">
                <span className="material-symbols-outlined text-white text-3xl">warning</span>
              </div>
              <h2 className="text-4xl font-black tracking-tighter uppercase leading-none">{t.simulation.realityCheck}</h2>
            </div>
            <div className="space-y-6 mb-10">
              <p className="text-xl font-bold tracking-tight uppercase">{t.simulation.realityCheckMsg.replace('{n}', sessionSpins)}</p>
              <p className="text-neutral-500 leading-relaxed italic border-l-4 border-primary pl-4">
                "{t.simulation.realityCheckDetail}"
              </p>
              <div className="bg-surface-container-low p-4">
                <p className="text-[10px] font-black uppercase tracking-widest text-neutral-400 mb-1">Estimated Reality Impact</p>
                <p className="text-2xl font-black text-error">-{formatValue(Math.abs(totalLoss + 12450))}</p>
              </div>
            </div>
            <div className="flex flex-col gap-4">
              <button 
                onClick={() => setShowRealityCheck(false)}
                className="bg-primary text-on-primary font-black py-4 uppercase tracking-widest hover:opacity-90 transition-opacity"
              >
                {t.simulation.realityCheckUnderstand}
              </button>
              <Link to="/intervention" className="text-center border-2 border-primary py-4 font-black uppercase tracking-widest hover:bg-primary hover:text-white transition-all">
                {t.simulation.realityCheckConsult}
              </Link>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

