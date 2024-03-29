# Bus factor calculator

Simple utility to calculate the bus-factor with

 * [AVL bus factor](https://arxiv.org/pdf/1604.06766.pdf) (based on file ownership)
 * [Pony factor](https://humbedooh.com/Chapter%203,%20part%20one_%20Codebase%20development%20resilience.pdf) (number of devs cover the >50% contributions, all time)
 * Dev power (number of devs as the factor of the biggest contributor)

Usage: execute application from a git repository.

Author mapping: You can make a `.git/bus-factor-alias` file which contains `alias@email.com,author@email.com` lines as mapping. (In case some of the authors use multiple emails). The first email (author) will be replaced with the second one.
