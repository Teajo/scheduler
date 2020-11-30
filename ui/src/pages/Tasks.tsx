import 'date-fns';
import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TablePagination from '@material-ui/core/TablePagination';
import TableRow from '@material-ui/core/TableRow';
import axios from 'axios';
import Picker from './TimeFilter';
import Checkbox from '@material-ui/core/Checkbox';
import { Button } from '@material-ui/core';

interface Column {
  id: string;
  label: string;
  minWidth?: number;
  align?: 'right';
  format?: (value: number) => string;
}

const columns: Column[] = [
  { id: 'id', label: 'ID', minWidth: 170 },
  { id: 'date', label: 'Date', minWidth: 100, format: (l) => new Date(l).toLocaleString() },
  { id: 'done', label: 'Done', minWidth: 100, format: (l) => String(l) },
];

const useStyles = makeStyles({
  root: {
    width: '100%',
  },
  container: {
    maxHeight: 620,
  },
});

export default function StickyHeadTable() {
  const classes = useStyles();
  const [page, setPage] = React.useState(0);
  const [rowsPerPage, setRowsPerPage] = React.useState(10);
  const [selectedStartDate, setSelectedStartDate] = React.useState<Date>(new Date());
  const [selectedEndDate, setSelectedEndDate] = React.useState<Date>(new Date());
  const [lastUpdate, setLastUpdate] = React.useState(new Date());
  const [tasks, setTasks] = React.useState([]);
  const [selected, setSelected] = React.useState<string[]>([]);
  const isSelected = (name: string) => selected.indexOf(name) !== -1;

  React.useEffect(() => {
    axios.get(`http://127.0.0.1:3000/tasks?endDate=${selectedEndDate.toISOString()}&startDate=${selectedStartDate.toISOString()}`)
      .then((res) => {
        const { data } = res.data;
        setTasks(data);
      })
      .catch((err) => {
        console.log(err);
      });
  }, [selectedStartDate, selectedEndDate, lastUpdate]);

  const handleChangePage = (event: any, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(+event.target.value);
    setPage(0);
  };

  const onStartChanged = (date: Date | null) => {
    if (date) {
      setSelectedStartDate(date);
    }
  };

  const onEndChanged = (date: Date | null) => {
    if (date) {
      setSelectedEndDate(date);
    }
  };

  const handleRowSelection = (id: string) => {
    if (!isSelected(id)) {
      setSelected([id, ...selected]);
    } else {
      const s = [...selected];
      const i = s.findIndex(t => t === id);
      s.splice(i, 1);
      setSelected(s);
    }
  };

  const onRemoveTask = async (id: string) => {
    try {
      const res = await axios.delete(`http://127.0.0.1:3000/tasks/${id}`);
    } catch (error) {
      console.log(error);
    } finally {
      setLastUpdate(new Date());
    }
  };

  return (
    <>
      <div style={{ textAlign: 'left' }}>
        <Picker label="Start" date={selectedStartDate} onChange={onStartChanged} />
        &nbsp;&nbsp;&nbsp;
        <Picker label="End" date={selectedEndDate} onChange={onEndChanged} />
      </div>

      <br />

      <Paper className={classes.root}>
        <TableContainer className={classes.container}>
          <Table stickyHeader aria-label="sticky table">
            <TableHead>
              <TableRow>
                <TableCell
                  style={{ minWidth: '170' }}
                />
                {columns.map((column) => (
                  <TableCell
                    key={column.id}
                    align={column.align}
                    style={{ minWidth: column.minWidth }}
                  >
                    {column.label}
                  </TableCell>
                ))}
                <TableCell style={{ minWidth: '170' }} />
              </TableRow>
            </TableHead>
            <TableBody>
              {tasks.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage).map((row: any) => (
                <>
                  <TableRow hover role="checkbox" tabIndex={-1} key={row.code}>
                    <TableCell padding="checkbox">
                      <Checkbox
                        checked={isSelected(row.id)}
                        onClick={() => handleRowSelection(row.id)}
                        inputProps={{ 'aria-label': 'select all desserts' }}
                      />
                    </TableCell>
                    {columns.map((column) => {
                      const value = row[column.id];
                      return (
                        <TableCell key={column.id} align={column.align}>
                          {column.format ? column.format(value) : value}
                        </TableCell>
                      );
                    })}
                    <TableCell>
                      <Button onClick={() => onRemoveTask(row.id)}>Remove</Button>
                    </TableCell>
                  </TableRow>
                </>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
        <TablePagination
          rowsPerPageOptions={[10, 25, 100]}
          component="div"
          count={tasks.length}
          rowsPerPage={rowsPerPage}
          page={page}
          onChangePage={handleChangePage}
          onChangeRowsPerPage={handleChangeRowsPerPage}
        />
      </Paper>
    </>
  );
}
