WikiRacer
=========

High level architecture
-----------------------
*Program does a bi directional BFS from the both start and end page concurrently.
*Uses the links API of media wiki for fetching all the links in a page. This is used to search forward from the source page.
*Uses the linkshere API of the media wiki for fetching all the links to a page. This is used to search backward from the destination page
*The forward and backward searches are done in parallel using go routines
*The results from the go routine are sent in channel to the parent
*Keeps the pages explored from the source in the left frontier, vice versa for the destination.
*When a page, which is in right frontier is reached from the forward search, we have found a path. Vice versa for the left frontier
*BFS is performed again from the source page to the left frontier edge page and right frontier edge page to the destination and get the path from source to destination

How to run
----------

Time spent
---------- 