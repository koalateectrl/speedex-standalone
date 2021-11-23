# speedex-standalone
Implementing SPEEDEX price computation engine in Golang as a standalone binary that exchanges can call.

Notes from Geoff About Tatonnement:

SCS Lab Implementation: https://github.com/scslab/speedex/blob/master/price_computation/

Stellar Implementation: https://github.com/gramseyer/stellar-core/blob/master/src/speedex/TatonnementOracle.cpp

Simple Simplex LP Solver: https://github.com/gramseyer/stellar-core/tree/master/src/simplex

Tatonnement

1) Collect DB of transactions for a given block
    1a) Create one vector for each group <sell A, buy B> (which is a different vector than <sell B, buy A>) as transactions come in
    1b) As we are building these vectors, keep each one in sorted order from minimum price.
    1c) Also keep track of cumulative amount up to that price.
2) Start with a price (ideal is integer - maybe 64 bit or less for edge cases and overflow and speed), but MVP can use floating point
3) Go through these vectors using binary search O(log n) to find how much would be bought/sold for each asset.
4) Update prices.
    4a) New price = Old price * (1 + Old price * asset quantity * (demand - supply) * STEP_SIZE)
    4b) Why weight by old price * asset quantity? It's to make sure when things are denominated differently they still take similar size steps. IE 1 dollar vs 100 pennies.
    4c) How to set STEP_SIZE? 
        4ci) Pick a small constant
        4cii) Run multiple Tatonnement in parallel with different constants
        4ciii) Use a heurtistic (such as L2 Norm or Infinity Norm) of Net Demand Vector (demand - supply) for Old price and New price and see if the norm decreases. If it increases then use a smaller STEP_SIZE.
5) Stopping criteria
    5a) Fixed number of iterations
    5b) Net Demand Norm less than some fixed amount.
6) LP Solver

Notes:
1) When collecting transactions from multiple threads, create this sorted DB (sell A, buy B, price, cumulative, metadata) per thread. Then do one MergeSort to combine them all together to get the entire transaction set.
2) Don't want demand to be too cliff-y, so smooth it out using a parameter alpha.
    2a) Example: alpha = 0.01. Then if limit price is 99% of current price then sell in full. If limit price is 99.5% of current price sell half of it. If limit price equals current price then sell none.
    2b) Will need to have 2 cumulative values in vector: cumulative amount and cumulative amoount * limit price offer



To Do:
    LP Solver after Tatonnement
