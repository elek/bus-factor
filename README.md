# Bus factor calculator

Simple utility to calculate the bus-factor with

## Metrics: 

### AVL Bus factor

See: [AVL bus factor](https://arxiv.org/pdf/1604.06766.pdf) (based on file ownership)

### Pony factor

Number of developers with biggest contributors, who created the >50% contributions

See: https://humbedooh.com/Chapter%203,%20part%20one_%20Codebase%20development%20resilience.pdf

### Dev power 

Number of devs as the factor of the biggest contributor.

Imagine that we have 5 contributors with 100, 35, 25, 20, 20 commits.

First, calculate the percentage for all contributions:

[ 50 %, 17.5 %, 12.5 %, 10 %, 10 % ]

First calculate the significance factor for each contributor. First contributor has weight 1, second is deducted with the difference of the contribution size (ratio).

[ 1, 1 * 35/100, 1 * 35/100 * 25/35, 1 * 35/100 * 25/35 * 20/25, 1 * 35/100 * 25/35 * 20/25 * 20/20 ]

Use these weights to adjust the percentage.


## Usage

Usage: execute application from a git repository.

Author mapping: You can make a `.git/bus-factor-alias` file which contains `alias@email.com,author@email.com` lines as mapping. (In case some of the authors use multiple emails). The first email (author) will be replaced with the second one.

