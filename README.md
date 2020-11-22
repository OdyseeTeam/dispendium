# dispendium
Receives api requests to send LBC to specific addresses. Used for the scalability of the LBRY Inc rewards program. 

Each cluster running with dispending contains a set of lbrycrd instances which 
are then used to round robin the sending across the instances associated with
the dispendium cluster. This allows for scalable, concurrent sending of
payments. 