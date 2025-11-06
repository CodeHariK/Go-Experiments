#include "position.h"
#include "tables.h"
#include "types.h"

int main()
{
    initialise_all_databases();
    zobrist::initialise_zobrist_keys();

    Position p;
    Position::set("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -", p);
    std::cout << p;

    MoveList<WHITE> list(p);

    for (Move m : list)
    {
        std::cout << m << "\n";
    }

    return 0;
}